package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupRoleRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	handler := handlers.NewRoleHandler(erpContext)
	group := r.Group("/role")
	group.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		group.GET("/permissions", handler.GetAllPermissionsHandler)
		group.GET("/list", handler.GetRolesHandler)
		group.DELETE("/:id", handler.DeleteRoleHandler)
		group.POST("/create", handler.CreateRoleHandler)
		group.PUT("/:id", handler.UpdateRoleHandler)
	}

}
