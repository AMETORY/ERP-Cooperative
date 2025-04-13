package middlewares

import (
	"ametory-cooperative/config"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

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

		c.Next()
	}
}

func ClosingBookMiddleware(ctx *context.ERPContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var start, end *time.Time
		if c.Request.Header.Get("start-date") != "" {
			startDate, err := time.Parse(time.RFC3339, c.Request.Header.Get("start-date"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
				c.Abort()
				return
			}
			start = &startDate
			c.Set("start-date", start)
		}
		if c.Request.Header.Get("end-date") != "" {
			endDate, err := time.Parse(time.RFC3339, c.Request.Header.Get("end-date"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
				c.Abort()
				return
			}
			end = &endDate
			c.Set("end-date", end)
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
