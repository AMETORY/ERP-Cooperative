package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupStockMovementRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	stockMovementHandler := handlers.NewStockMovementHandler(erpContext)

	stockMovementGroup := r.Group("/stock-movement")
	stockMovementGroup.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		stockMovementGroup.POST("/create", stockMovementHandler.CreateStockMovementHandler)
		stockMovementGroup.GET("/list", stockMovementHandler.GetStockMovementHandler)
		stockMovementGroup.GET("/:id", stockMovementHandler.GetStockMovementByIdHandler)
		stockMovementGroup.PUT("/:id", stockMovementHandler.UpdateStockMovementHandler)
		stockMovementGroup.DELETE("/:id", stockMovementHandler.DeleteStockMovementHandler)
	}
}
