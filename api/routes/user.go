package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	handler := handlers.NewUserHandler(erpContext)
	group := r.Group("/user")
	group.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		group.POST("/activity", handler.CreateActivityHandler)
		group.POST("/activity/:id/finish", handler.FinishActivityHandler)
		group.POST("/clock-in", handler.ClockInHandler)
		group.POST("/last-clock-in", handler.GetLastClockinHandler)
		group.POST("/clock-out", handler.ClockOutHandler)
		group.GET("/activities", handler.GetActivityHandler)
		group.POST("/break", handler.BreakHandler)
		group.POST("/break-off", handler.BreakOffHandler)
		group.GET("/:id", handler.GetUserDetailHandler)
		group.GET("/list", handler.GetUserListHandler)
	}
}
