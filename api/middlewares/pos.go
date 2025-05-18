package middlewares

import (
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func PosMiddleware(ctx *context.ERPContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
