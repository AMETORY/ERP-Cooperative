package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupProductAttributeRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	productAttributeHandler := handlers.NewProductAttributeHandler(erpContext)
	productAttributeGroup := r.Group("/product-attribute")
	productAttributeGroup.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		productAttributeGroup.GET("/list", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:product_attribute:read"}), productAttributeHandler.GetProductAttributeHandler)
		productAttributeGroup.GET("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:product_attribute:read"}), productAttributeHandler.GetProductAttributeByIdHandler)
		productAttributeGroup.POST("/create", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:product_attribute:create"}), productAttributeHandler.CreateProductAttributeHandler)
		productAttributeGroup.PUT("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:product_attribute:update"}), productAttributeHandler.UpdateProductAttributeHandler)
		productAttributeGroup.DELETE("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:product_attribute:delete"}), productAttributeHandler.DeleteProductAttributeHandler)
	}
}
