package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupMerchantRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	merchantHandler := handlers.NewMerchantHandler(erpContext)
	merchantGroup := r.Group("/merchant")
	merchantGroup.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		merchantGroup.GET("/list", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:read"}), merchantHandler.ListMerchantsHandler)
		merchantGroup.GET("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:read"}), merchantHandler.GetMerchantHandler)
		merchantGroup.GET("/:id/products", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:read"}), merchantHandler.GetMerchantProductsHandler)
		merchantGroup.POST("/create", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:create"}), merchantHandler.CreateMerchantHandler)
		merchantGroup.PUT("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:update"}), merchantHandler.UpdateMerchantHandler)
		merchantGroup.PUT("/:id/add-product", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:update"}), merchantHandler.AddProductMerchantHandler)
		merchantGroup.DELETE("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:delete"}), merchantHandler.DeleteMerchantHandler)
		merchantGroup.DELETE("/:id/delete-product", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:delete"}), merchantHandler.DeleteProductsFromMerchantHandler)
	}
}
