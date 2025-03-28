package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	authHandler := handlers.NewAuthHandler(erpContext)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authHandler.LoginHandler)
		authGroup.POST("/register", authHandler.RegisterHandler)
		authGroup.GET("/verification/:token", authHandler.VerificationEmailHandler)
		authGroup.POST("/change-password", middlewares.AuthMiddleware(erpContext, false), authHandler.ChangePasswordHandler)
		authGroup.GET("/profile", middlewares.AuthMiddleware(erpContext, false), authHandler.GetProfile)
		authGroup.PUT("/profile", middlewares.AuthMiddleware(erpContext, false), authHandler.UpdateProfile)
	}
}
