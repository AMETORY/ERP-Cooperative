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
		salesGroup.DELETE("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"sales:sales:delete"}), salesHandler.DeleteSalesHandler)
	}
}
