package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupStockOpnameRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	stockOpnameHandler := handlers.NewStockOpnameHandler(erpContext)
	stockOpnameGroup := r.Group("/stock-opname")
	stockOpnameGroup.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		stockOpnameGroup.GET("/list", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:stock_opname:read"}), stockOpnameHandler.ListStockOpnamesHandler)
		stockOpnameGroup.GET("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:stock_opname:read"}), stockOpnameHandler.GetStockOpnameHandler)
		stockOpnameGroup.POST("/create", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:stock_opname:create"}), stockOpnameHandler.CreateStockOpnameHandler)
		stockOpnameGroup.PUT("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:stock_opname:update"}), stockOpnameHandler.UpdateStockOpnameHandler)

		stockOpnameGroup.PUT("/:id/add-item", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:stock_opname:update"}), stockOpnameHandler.AddItemHandler)
		stockOpnameGroup.PUT("/:id/update-item/:detailId", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:stock_opname:update"}), stockOpnameHandler.UpdateItemHandler)
		stockOpnameGroup.DELETE("/:id/delete-item/:detailId", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:stock_opname:update"}), stockOpnameHandler.DeleteItemHandler)
		stockOpnameGroup.PUT("/:id/complete", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:stock_opname:update"}), stockOpnameHandler.CompleteStockOpnameHandler)
		stockOpnameGroup.DELETE("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:stock_opname:delete"}), stockOpnameHandler.DeleteStockOpnameHandler)
	}
}
