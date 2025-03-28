package middlewares

import (
	"ametory-cooperative/config"
	"net/http"
	"strings"

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
		ctx.DB.Preload("Company").Where("connected_to = ? and company_id = ?", token.Claims.(*jwt.StandardClaims).Id, c.Request.Header.Get("ID-Company")).Find(&member)
		c.Set("user", user)
		c.Set("member", member)
		c.Set("memberID", member.ID)
		c.Next()
	}
}
