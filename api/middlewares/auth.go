package middlewares

import (
	"ametory-cooperative/app_models"
	"ametory-cooperative/config"
	"ametory-cooperative/services"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"github.com/golang-jwt/jwt"
)

const (
	UsageLimit       = 10000
	UsageWindow      = 6 * time.Hour
	CooldownDuration = 4 * time.Hour
)

var exceptionPaths = []string{
	"/api/v1/auth/profile",
	"/api/v1/setting",
	"/api/v1/create/company",
}

func AuthMiddleware(ctx *context.ERPContext, checkCompany bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		splitToken := strings.Split(authHeader, "Bearer ")

		if len(splitToken) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}
		reqToken := splitToken[1]
		if reqToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Token unsplited"})
			c.Abort()
			return
		}
		// fmt.Println("reqToken: ", reqToken)

		token, err := jwt.ParseWithClaims(reqToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.App.Server.SecretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if c.Request.Header.Get("ID-Company") == "" && checkCompany {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID is required"})
			c.Abort()
		}

		c.Set("companyID", c.Request.Header.Get("ID-Company"))
		c.Set("userID", token.Claims.(*jwt.StandardClaims).Id)
		user := models.UserModel{}
		ctx.DB.Find(&user, "id = ?", token.Claims.(*jwt.StandardClaims).Id)
		var member models.CooperativeMemberModel
		ctx.DB.Preload("Company").Where("connected_to = ? and company_id = ?", token.Claims.(*jwt.StandardClaims).Id, c.Request.Header.Get("ID-Company")).Preload("Role.Permissions").Find(&member)
		c.Set("user", user)
		c.Set("member", member)
		c.Set("memberID", member.ID)

		for _, path := range exceptionPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		if !config.App.Server.UsedBasedQuota {
			c.Next()
			return
		}

		company := app_models.CustomSettingModel{}
		ctx.DB.Find(&company, "id = ?", c.Request.Header.Get("ID-Company"))

		if company.IsPremium && company.PremiumExpiredAt != nil {
			if company.PremiumExpiredAt.After(time.Now()) {
				c.Next()
				return
			}

		}
		err = UsageRateLimitMiddleware(ctx, token.Claims.(*jwt.StandardClaims).Id)
		if err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}

func UsageRateLimitMiddleware(ctx *context.ERPContext, userID string) error {
	if userID == "" {
		return errors.New("missing user ID")
	}

	usageKey := fmt.Sprintf("user:%s:usage", userID)
	cooldownKey := fmt.Sprintf("user:%s:cooldown", userID)

	// Cek apakah sedang cooldown
	cooling, err := services.REDIS.Exists(*ctx.Ctx, cooldownKey).Result()
	if err != nil {
		return errors.New("redis error")
	}
	if cooling == 1 {

		cooldownTTL, err := services.REDIS.TTL(*ctx.Ctx, cooldownKey).Result()
		if err != nil {
			return errors.New("redis error")
		}
		return fmt.Errorf("rate limit exceeded. Please wait for cooldown, remain time: %v", time.Duration(cooldownTTL))
	}

	// Ambil usage count
	count, err := services.REDIS.Get(*ctx.Ctx, usageKey).Int()
	if err == redis.Nil {
		// Belum ada, set ke 1 dengan TTL 6 jam
		services.REDIS.Set(*ctx.Ctx, usageKey, 1, UsageWindow)
	} else if err != nil {
		return errors.New("redis error")
	} else if count >= UsageLimit {
		// Set cooldown 4 jam
		services.REDIS.Set(*ctx.Ctx, cooldownKey, true, CooldownDuration)
		return errors.New("rate limit exceeded. Cooldown started")
	} else {
		// Tambah hit
		services.REDIS.Incr(*ctx.Ctx, usageKey)
	}

	return nil
}

func ClosingBookMiddleware(ctx *context.ERPContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var start, end *time.Time
		if c.Request.Header.Get("start-date") != "" {
			startDate, err := time.Parse(time.RFC3339, c.Request.Header.Get("start-date"))
			if err == nil {
				start = &startDate
				c.Set("start-date", start)
			}

		}
		if c.Request.Header.Get("end-date") != "" {
			endDate, err := time.Parse(time.RFC3339, c.Request.Header.Get("end-date"))
			if err == nil {
				end = &endDate
				c.Set("end-date", end)
			}

		}
		fmt.Println("CHECK PERIODE #2", start, end)
		if start == nil || end == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Accounting period is not set. Please specify the start and end dates in the request headers."})
			c.Abort()
			return
		}
		var closingBook models.ClosingBook
		ctx.DB.Where("start_date  >= ? and end_date <= ?", start, end).First(&closingBook)
		c.Set("closingBook", &closingBook)
		c.Next()
	}
}
