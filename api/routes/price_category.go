package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupPriceCategoryRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	priceCategoryHandler := handlers.NewPriceCategoryHandler(erpContext)
	priceCategoryGroup := r.Group("/price-category")
	priceCategoryGroup.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		priceCategoryGroup.GET("/list", priceCategoryHandler.ListPriceCategoriesHandler)
		priceCategoryGroup.GET("/:id", priceCategoryHandler.GetPriceCategoryHandler)
		priceCategoryGroup.POST("/create", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:price_category:create"}), priceCategoryHandler.CreatePriceCategoryHandler)
		priceCategoryGroup.PUT("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:price_category:update"}), priceCategoryHandler.UpdatePriceCategoryHandler)
		priceCategoryGroup.DELETE("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:price_category:delete"}), priceCategoryHandler.DeletePriceCategoryHandler)
	}
}
