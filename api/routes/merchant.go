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
		merchantGroup.GET("/:id/users", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:read"}), merchantHandler.GetMerchantUsersHandler)
		merchantGroup.PUT("/:id/add-user", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:update"}), merchantHandler.AddUserMerchantHandler)
		merchantGroup.DELETE("/:id/delete-user", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:delete"}), merchantHandler.DeleteUserFromMerchantHandler)
		merchantGroup.PUT("/:id/add-desk", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:update"}), merchantHandler.AddDeskMerchantHandler)
		merchantGroup.GET("/:id/desk", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:update"}), merchantHandler.GetDeskMerchantHandler)
		merchantGroup.PUT("/:id/desk/:deskId", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:update"}), merchantHandler.UpdateDeskMerchantHandler)
		merchantGroup.DELETE("/:id/desk/:deskId", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:update"}), merchantHandler.DeleteDeskMerchantHandler)
		merchantGroup.GET("/:id/layouts", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:update"}), merchantHandler.GetLayoutsMerchantHandler)
		merchantGroup.GET("/:id/layout/:layoutId", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:update"}), merchantHandler.GetLayoutDetailMerchantHandler)
		merchantGroup.PUT("/:id/add-layout", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:update"}), merchantHandler.AddLayoutMerchantHandler)
		merchantGroup.PUT("/:id/update-layout/:layoutId", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:update"}), merchantHandler.UpdateLayoutMerchantHandler)
		merchantGroup.DELETE("/:id/delete-layout/:layoutId", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:merchant:update"}), merchantHandler.DeleteLayoutMerchantHandler)
	}
}
