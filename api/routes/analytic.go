package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetUpAnalyticRoutes(r *gin.RouterGroup, ctx *context.ERPContext) {
	handler := handlers.NewAnalyticHandler(ctx)
	group := r.Group("/analytic")
	group.Use(middlewares.AuthMiddleware(ctx, false))
	{
		group.GET("/popular-product", handler.PopularProductHandler)
		group.GET("/monthly-sales", handler.GetMonthlySalesReportHandler)
		group.GET("/monthly-purchase", handler.GetMonthlyPurchaseReportHandler)
		group.GET("/monthly-sales-purchase", handler.GetMonthlySalesPurchaseReportHandler)
		group.GET("/weekly-sales", handler.GetWeeklySalesReportHandler)
		group.GET("/weekly-purchase", handler.GetWeeklyPurchaseReportHandler)
		group.GET("/weekly-sales-purchase", handler.GetWeeklySalesPurchaseReportHandler)
		group.GET("/sales-time-range", handler.GetSalesTimeRangeHandler)
		group.GET("/purchase-time-range", handler.GetPurchaseTimeRangeHandler)
	}
}
