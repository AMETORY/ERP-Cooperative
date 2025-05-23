package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupReportRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	reportHandler := handlers.NewReportHandler(erpContext)
	reportGroup := r.Group("/report")
	reportGroup.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		reportGroup.POST("/cogs", reportHandler.GetCogsHandler)
		reportGroup.POST("/profit-loss", reportHandler.GetProfitLossHandler)
		reportGroup.POST("/balance-sheet", reportHandler.GetBalanceSheetHandler)
		reportGroup.POST("/capital-change", reportHandler.CapitalChangeHandler)
		reportGroup.POST("/cash-flow", reportHandler.CashFlowHandler)
		reportGroup.POST("/trial-balance", reportHandler.TrialBalanceHandler)
		reportGroup.GET("/closing-book", reportHandler.GetClosingBooksHandler)
		reportGroup.POST("/closing-book", reportHandler.CreateClosingBookHandler)
		reportGroup.GET("/closing-book/:id", reportHandler.GetClosingBookDetailHandler)
		reportGroup.DELETE("/closing-book/:id", reportHandler.DeleteClosingBooklHandler)
		reportGroup.PUT("/closing-book/:id/generate", reportHandler.GenerateClosingBookHandler)
		reportGroup.POST("/product-sales-customers", reportHandler.GetProductSalesCustomersHandler)
		reportGroup.POST("/account-receivable-ledger", reportHandler.GetAccountReceivableLedgerHandler)
		reportGroup.POST("/account-payable-ledger", reportHandler.GetAccountPayabledgerHandler)
	}
}
