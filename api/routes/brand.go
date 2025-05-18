package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupBrandRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	brandHandler := handlers.NewBrandHandler(erpContext)
	brandGroup := r.Group("/brand")
	brandGroup.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		brandGroup.GET("/list", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:brand:read"}), brandHandler.GetBrandHandler)
		brandGroup.GET("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:brand:read"}), brandHandler.GetBrandHandler)
		brandGroup.POST("/create", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:brand:create"}), brandHandler.CreateBrandHandler)
		brandGroup.PUT("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:brand:update"}), brandHandler.UpdateBrandHandler)
		brandGroup.DELETE("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:brand:delete"}), brandHandler.DeleteBrandHandler)
	}
}
