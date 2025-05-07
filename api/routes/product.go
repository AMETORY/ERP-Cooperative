package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	productHandler := handlers.NewProductHandler(erpContext)
	productGroup := r.Group("/product")
	productGroup.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		productGroup.GET("/list", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:product:read"}), productHandler.ListProductsHandler)
		productGroup.GET("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:product:read"}), productHandler.GetProductHandler)
		productGroup.GET("/:id/variant", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:product:read"}), productHandler.GetProductVariantHandler)
		productGroup.POST("/:id/variant", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:product:read"}), productHandler.CreateProductVariantHandler)
		productGroup.PUT("/:id/unit", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:product:update"}), productHandler.AddProductUnitHandler)
		productGroup.DELETE("/:id/unit/:unitId", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:product:update"}), productHandler.DeleteProductUnitHandler)
		productGroup.GET("/:id/discount", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:product:read"}), productHandler.GetProductDiscountHandler)
		productGroup.POST("/:id/discount", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:product:update"}), productHandler.AddDiscountHandler)
		productGroup.POST("/create", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:product:create"}), productHandler.CreateProductHandler)
		productGroup.PUT("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:product:update"}), productHandler.UpdateProductHandler)
		productGroup.PUT("/:id/price", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:product:update"}), productHandler.AddPriceProductHandler)
		productGroup.DELETE("/:id/price/:priceId", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:product:update"}), productHandler.DeletePriceProductHandler)
		productGroup.DELETE("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:product:delete"}), productHandler.DeleteProductHandler)
	}
}
