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
		transactionGroup.GET("/list", middlewares.RbacUserMiddleware(erpContext, false, []string{"finance:transaction:read"}), transactionHandler.ListTransactions)
		transactionGroup.GET("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"finance:transaction:read"}), transactionHandler.GetTransaction)
		transactionGroup.POST("/create", middlewares.ClosingBookMiddleware(erpContext), middlewares.RbacUserMiddleware(erpContext, false, []string{"finance:transaction:create"}), transactionHandler.CreateTransaction)
		transactionGroup.PUT("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"finance:transaction:update"}), transactionHandler.UpdateTransaction)
		transactionGroup.DELETE("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"finance:transaction:delete"}), transactionHandler.DeleteTransaction)
	}
}
