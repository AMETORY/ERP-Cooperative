package routes

import (
	"ametory-cooperative/api/handlers/pos"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupPosRoutes(router *gin.RouterGroup, erpContext *context.ERPContext) {
	handler := pos.NewPosHandler(erpContext)
	group := router.Group("/pos")
	group.Use(middlewares.AuthMiddleware(erpContext, false), middlewares.PosMiddleware(erpContext))
	{
		group.GET("/merchants", handler.GetMerchantsHandler)
		group.GET("/merchant/:id/orders", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:pos:cashier"}), handler.GetOrdersHandler)
		group.GET("/merchant/:id/stations", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:pos:cashier"}), handler.GetStationsHandler)
		group.GET("/merchant/:id/station/:stationId", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:pos:cashier"}), handler.GetStationDetailHandler)
		group.GET("/merchant/:id/station/:stationId/orders", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:pos:cashier"}), handler.GetStationOrdersHandler)
		group.POST("/merchant/:id/order", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:pos:cashier"}), handler.CreateOrderHandler)
		group.GET("/merchant/:id/layouts", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:pos:cashier"}), handler.GetLayoutsMerchantHandler)
		group.GET("/merchant/:id/layout/:layoutId", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:pos:cashier"}), handler.GetLayoutDetailMerchantHandler)
		group.PUT("/merchant/:id/table/:tableId/update-status", middlewares.RbacUserMiddleware(erpContext, false, []string{"order:pos:cashier"}), handler.UpdateStatusTableHandler)
	}
}
