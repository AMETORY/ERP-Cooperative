package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupTaxRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	taxHandler := handlers.NewTaxHandler(erpContext)
	taxGroup := r.Group("/tax")
	taxGroup.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		taxGroup.GET("/list", middlewares.RbacUserMiddleware(erpContext, false, []string{"finance:tax:read"}), taxHandler.GetTaxHandler)
		taxGroup.GET("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"finance:tax:read"}), taxHandler.GetTaxByIdHandler)
		taxGroup.POST("/create", middlewares.RbacUserMiddleware(erpContext, false, []string{"finance:tax:create"}), taxHandler.CreateTaxHandler)
		taxGroup.PUT("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"finance:tax:update"}), taxHandler.UpdateTaxHandler)
		taxGroup.DELETE("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"finance:tax:delete"}), taxHandler.DeleteTaxHandler)
	}
}
