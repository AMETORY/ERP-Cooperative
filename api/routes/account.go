package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupAccountRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	accountHandler := handlers.NewAccountHandler(erpContext)
	accountGroup := r.Group("/account")
	accountGroup.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		accountGroup.GET("/chart-of-accounts", accountHandler.GetChartOfAccounts)
		accountGroup.POST("/create", accountHandler.CreateAccountHandler)
		accountGroup.GET("/list", accountHandler.GetAccountHandler)
		accountGroup.GET("/:id", accountHandler.GetAccountByIdHandler)
		accountGroup.PUT("/:id", accountHandler.UpdateAccountHandler)
		accountGroup.DELETE("/:id", accountHandler.DeleteAccountHandler)

	}
}
