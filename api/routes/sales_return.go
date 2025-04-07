package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupSalesReturnRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	handler := handlers.NewSalesReturnHandler(erpContext)
	group := r.Group("/sales-return")
	group.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		group.GET("/list", handler.GetSalesReturnListHandler)
		group.GET("/:id", handler.GetSalesReturnHandler)
		group.POST("/create", handler.CreateSalesReturnHandler)
		group.PUT("/:id", handler.UpdateSalesReturnHandler)
		group.PUT("/:id/release", handler.ReleaseSalesReturnHandler)
		group.PUT("/:id/add-item", handler.AddItemSalesReturnHandler)
		group.PUT("/:id/update-item/:itemId", handler.UpdateItemSalesReturnHandler)
		group.PUT("/:id/delete-item/:itemId", handler.DeleteItemSalesReturnHandler)
		group.DELETE("/:id", handler.DeleteSalesReturnHandler)

	}
}
