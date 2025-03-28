package middlewares

import (
	"errors"

	"github.com/AMETORY/ametory-erp-modules/auth"
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RbacSuperAdminMiddleware(erpContext *context.ERPContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		ok, _ := checkSuperAdmin(erpContext, userID.(string), c.GetHeader("ID-Company"))
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

func checkSuperAdmin(erpContext *context.ERPContext, userID string, companyID string) (bool, error) {
	var admin models.UserModel

	// Cari pengguna beserta peran dan izin
	if err := erpContext.DB.Preload("Roles", func(db *gorm.DB) *gorm.DB {
		return db.Where("company_id = ?", companyID)
	}).First(&admin, "id = ?", userID).Error; err != nil {
		return false, errors.New("admin not found")
	}

	// Periksa izin
	for _, role := range admin.Roles {
		if role.IsSuperAdmin {
			return true, nil
		}
	}

	return false, nil
}
