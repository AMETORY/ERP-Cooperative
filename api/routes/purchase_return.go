package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupPurchaseReturnRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	handler := handlers.NewPurchaseReturnHandler(erpContext)
	group := r.Group("/purchase-return")
	group.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		group.GET("/list", handler.GetPurchaseReturnListHandler)
		group.GET("/:id", handler.GetPurchaseReturnHandler)
		group.POST("/create", handler.CreatePurchaseReturnHandler)
		group.PUT("/:id/add-item", handler.AddItemPurchaseReturnHandler)
		group.PUT("/:id/update-item/:itemId", handler.UpdateItemPurchaseReturnHandler)
		group.PUT("/:id/delete-item/:itemId", handler.DeleteItemPurchaseReturnHandler)
		group.DELETE("/:id", handler.DeletePurchaseReturnHandler)

	}
}
