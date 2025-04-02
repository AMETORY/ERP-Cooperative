package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupWarehouseRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	handler := handlers.NewWarehouseHandler(erpContext)
	warehouseGroup := r.Group("/warehouse")
	warehouseGroup.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		warehouseGroup.GET("/list", handler.ListWarehousesHandler)
		warehouseGroup.GET("/:id", handler.GetWarehouseByIdHandler)
		warehouseGroup.POST("/create", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:warehouse:create"}), handler.CreateWarehouseHandler)
		warehouseGroup.PUT("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:warehouse:update"}), handler.UpdateWarehouseHandler)
		warehouseGroup.DELETE("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:warehouse:delete"}), handler.DeleteWarehouseHandler)
	}
}
