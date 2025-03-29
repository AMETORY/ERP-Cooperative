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
		accountGroup.GET("/account-types", accountHandler.GetAccountTypesHandler)
		accountGroup.GET("/get-code", accountHandler.GetCodeHandler)
		accountGroup.POST("/create", accountHandler.CreateAccountHandler)
		accountGroup.GET("/list", accountHandler.GetAccountHandler)
		accountGroup.GET("/:id", accountHandler.GetAccountByIdHandler)
		accountGroup.GET("/:id/report", accountHandler.GetAccountReportHandler)
		accountGroup.PUT("/:id", accountHandler.UpdateAccountHandler)
		accountGroup.DELETE("/:id", accountHandler.DeleteAccountHandler)

	}
}
