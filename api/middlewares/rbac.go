package middlewares

import (
	"github.com/AMETORY/ametory-erp-modules/auth"
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func RbacSuperAdminMiddleware(erpContext *context.ERPContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		rbacSrv, ok := erpContext.RBACService.(*auth.RBACService)
		if !ok {
			c.JSON(500, gin.H{"message": "Auth service is not available"})
			c.Abort()
			return
		}
		ok, _ = rbacSrv.CheckSuperAdminPermission(userID.(string))
		if !ok {
			c.JSON(403, gin.H{"message": "Forbidden"})
			c.Abort()
		}
		c.Next()
	}
}
func RbacUserMiddleware(erpContext *context.ERPContext, isAdmin bool, permissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		rbacSrv, ok := erpContext.RBACService.(*auth.RBACService)
		if !ok {
			c.JSON(500, gin.H{"message": "Auth service is not available"})
			c.Abort()
			return
		}
		if isAdmin {
			ok, _ := rbacSrv.CheckAdminPermission(userID.(string), permissions)
			if !ok {
				c.JSON(403, gin.H{"message": "Forbidden"})
				c.Abort()
				return
			}
		} else {
			ok, _ := rbacSrv.CheckPermission(userID.(string), permissions)
			if !ok {
				c.JSON(403, gin.H{"message": "Forbidden"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
