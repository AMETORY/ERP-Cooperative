package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupJournalRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	journalHandler := handlers.NewJournalHandler(erpContext)
	journalGroup := r.Group("/journal")
	journalGroup.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		journalGroup.GET("/list", middlewares.RbacUserMiddleware(erpContext, false, []string{"finance:journal:read"}), journalHandler.ListJournalsHandler)
		journalGroup.GET("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"finance:journal:read"}), journalHandler.GetJournalHandler)
		journalGroup.POST("/create", middlewares.RbacUserMiddleware(erpContext, false, []string{"finance:journal:create"}), journalHandler.CreateJournalHandler)
		journalGroup.PUT("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"finance:journal:update"}), journalHandler.UpdateJournalHandler)
		journalGroup.PUT("/:id/add-transaction", middlewares.RbacUserMiddleware(erpContext, false, []string{"finance:journal:update"}), journalHandler.AddTransactionHandler)
		journalGroup.DELETE("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"finance:journal:delete"}), journalHandler.DeleteJournalHandler)
	}

}
