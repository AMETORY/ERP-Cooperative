package middlewares

import (
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

func PosMiddleware(ctx *context.ERPContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var merchant *models.MerchantModel
		if c.Request.Header.Get("ID-Merchant") != "" {
			var merch models.MerchantModel
			ctx.DB.Find(&merch, "id = ?", c.Request.Header.Get("ID-Merchant"))
			merchant = &merch
			c.Set("merchant", merchant)
			c.Set("merchantID", c.Request.Header.Get("ID-Merchant"))
		}
		c.Next()
	}
}
