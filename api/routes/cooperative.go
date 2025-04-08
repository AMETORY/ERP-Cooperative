package routes

import (
	cooperative_handler "ametory-cooperative/api/handlers/cooperative"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupCooperativeRoutes(r *gin.RouterGroup, ctx *context.ERPContext) {
	// Cooperative Members
	memberHandler := cooperative_handler.NewCooperativeMemberHandler(ctx)
	memberGroup := r.Group("/cooperative/member")
	memberGroup.Use(middlewares.AuthMiddleware(ctx, false))
	{
		memberGroup.GET("/list", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:cooperative_member:read"}), memberHandler.GetMembersHandler)
		memberGroup.GET("/:id", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:cooperative_member:read"}), memberHandler.GetMemberHandler)
		memberGroup.POST("/create", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:cooperative_member:create"}), memberHandler.CreateMemberHandler)
		memberGroup.PUT("/:id", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:cooperative_member:update"}), memberHandler.UpdateMemberHandler)
		memberGroup.DELETE("/:id", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:cooperative_member:delete"}), memberHandler.DeleteMemberHandler)
		// memberGroup.POST("/invite", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:cooperative_member:invite"}), memberHandler.InviteMemberHandler)
		// memberGroup.POST("/approve/:id", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:cooperative_member:approval"}), memberHandler.ApproveMemberHandler)
	}

}
