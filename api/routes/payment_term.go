package routes

import (
	"ametory-cooperative/api/handlers"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupPaymentTermRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	handler := handlers.NewPaymentTermHandler(erpContext)
	group := r.Group("/payment-term")
	{
		group.GET("/list", handler.GetPaymentTermsHandler)
		group.GET("/group", handler.GetPaymentTermsGroupHandler)
	}
}
