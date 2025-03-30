package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetContactRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	contactHandler := handlers.NewContactHandler(erpContext)

	contactGroup := r.Group("/contact")
	contactGroup.Use(middlewares.AuthMiddleware(erpContext, true))
	{
		contactGroup.GET("/list", middlewares.RbacUserMiddleware(erpContext, false, []string{"contact:all:read"}), contactHandler.GetContactsHandler)
		contactGroup.GET("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"contact:all:read"}), contactHandler.GetContactHandler)
		contactGroup.POST("/create", middlewares.RbacUserMiddleware(erpContext, false, []string{"contact:all:create"}), contactHandler.CreateContactHandler)
		contactGroup.PUT("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"contact:all:update"}), contactHandler.UpdateContactHandler)
		contactGroup.DELETE("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"contact:all:delete"}), contactHandler.GetContactsHandler)
	}
}
