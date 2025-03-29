package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupTransactionRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	transactionHandler := handlers.NewTransactionHandler(erpContext) // TODO: inject transactionHandler
	transactionGroup := r.Group("/transaction")
	transactionGroup.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		transactionGroup.GET("/list", transactionHandler.ListTransactions)
		transactionGroup.GET("/:id", transactionHandler.GetTransaction)
		transactionGroup.POST("/create", transactionHandler.CreateTransaction)
		transactionGroup.PUT("/:id", transactionHandler.UpdateTransaction)
		transactionGroup.DELETE("/:id", transactionHandler.DeleteTransaction)
	}
}
