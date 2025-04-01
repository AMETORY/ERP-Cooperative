package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupSalesRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	salesHandler := handlers.NewSalesHandler(erpContext)
	salesGroup := r.Group("/sales")
	salesGroup.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		salesGroup.GET("/list", middlewares.RbacUserMiddleware(erpContext, false, []string{"sales:sales:read"}), salesHandler.GetSalesHandler)
		salesGroup.GET("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"sales:sales:read"}), salesHandler.GetSalesByIdHandler)
		salesGroup.POST("/create", middlewares.RbacUserMiddleware(erpContext, false, []string{"sales:sales:create"}), salesHandler.CreateSalesHandler)
		salesGroup.PUT("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"sales:sales:update"}), salesHandler.UpdateSalesHandler)
		salesGroup.PUT("/:id/add-item", middlewares.RbacUserMiddleware(erpContext, false, []string{"sales:sales:update"}), salesHandler.AddItemHandler)
		salesGroup.PUT("/:id/update-item/:itemId", middlewares.RbacUserMiddleware(erpContext, false, []string{"sales:sales:update"}), salesHandler.UpdateItemHandler)
		salesGroup.GET("/:id/items", middlewares.RbacUserMiddleware(erpContext, false, []string{"sales:sales:read"}), salesHandler.GetItemsHandler)
		salesGroup.DELETE("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"sales:sales:delete"}), salesHandler.DeleteSalesHandler)
	}
}
