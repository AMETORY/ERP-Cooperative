package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupPurchaseRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	handler := handlers.NewPurchaseHandler(erpContext)
	group := r.Group("/purchase")
	group.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		group.GET("/list", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:purchase:read"}), handler.GetPurchaseHandler)
		group.GET("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:purchase:read"}), handler.GetPurchaseByIdHandler)
		group.POST("/create", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:purchase:create"}), handler.CreatePurchaseHandler)
		group.PUT("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:purchase:update"}), handler.UpdatePurchaseHandler)
		group.PUT("/:id/payment", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:purchase:update"}), handler.PaymentHandler)
		group.PUT("/:id/add-item", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:purchase:update"}), handler.AddItemHandler)
		group.PUT("/:id/post", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:purchase:update"}), handler.PostPurchaseHandler)
		group.PUT("/:id/publish", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:purchase:update"}), handler.PublishPurchaseHandler)
		group.DELETE("/:id/delete-item/:itemId", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:purchase:update"}), handler.DeleteItemHandler)
		group.PUT("/:id/update-item/:itemId", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:purchase:update"}), handler.UpdateItemHandler)
		group.GET("/:id/items", middlewares.RbacUserMiddleware(erpContext, false, []string{"inventory:purchase:read"}), handler.GetItemsHandler)
		group.DELETE("/delete/:id", handler.DeletePurchaseHandler)
	}
}
