package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	sales "ametory-cooperative/api/handlers/sales"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupSalesRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	salesHandler := handlers.NewSalesHandler(erpContext)
	salesDashboardHandler := sales.NewSalesDashboardHandler(erpContext)
	salesGroup := r.Group("/sales")
	salesGroup.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		salesGroup.GET("/list", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:sales:read"}), salesHandler.GetSalesHandler)
		salesGroup.GET("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:sales:read"}), salesHandler.GetSalesByIdHandler)
		salesGroup.GET("/:id/pdf", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:sales:read"}), salesHandler.DownloadPdfHandler)
		salesGroup.POST("/create", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:sales:create"}), salesHandler.CreateSalesHandler)
		salesGroup.PUT("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:sales:update"}), salesHandler.UpdateSalesHandler)
		salesGroup.PUT("/:id/add-item", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:sales:update"}), salesHandler.AddItemHandler)
		salesGroup.PUT("/:id/payment", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:sales:update"}), salesHandler.PaymentHandler)
		salesGroup.PUT("/:id/post", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:sales:update"}), salesHandler.PostInvoiceHandler)
		salesGroup.PUT("/:id/publish", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:sales:update"}), salesHandler.PublishSalesHandler)
		salesGroup.DELETE("/:id/delete-item/:itemId", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:sales:update"}), salesHandler.DeleteItemHandler)
		salesGroup.PUT("/:id/update-item/:itemId", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:sales:update"}), salesHandler.UpdateItemHandler)
		salesGroup.GET("/:id/items", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:sales:read"}), salesHandler.GetItemsHandler)
		salesGroup.DELETE("/:id", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:sales:delete"}), salesHandler.DeleteSalesHandler)

		dashboardGroup := salesGroup.Group("/dashboard")
		dashboardGroup.Use(middlewares.AuthMiddleware(erpContext, false))
		{
			dashboardGroup.POST("/summary", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:sales:read"}), salesDashboardHandler.GetDashboardSummaryHandler)
		}
	}
}
