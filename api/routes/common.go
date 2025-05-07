package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupCommonRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	commonHandler := handlers.NewCommonHandler(erpContext)
	{
		r.POST("/file/upload", middlewares.AuthMiddleware(erpContext, false), commonHandler.UploadFileHandler)
		r.POST("/file/upload-base64", middlewares.AuthMiddleware(erpContext, false), commonHandler.UploadFileFromBase64Handler)
		r.GET("/company/users", middlewares.AuthMiddleware(erpContext, false), commonHandler.GetCompanyUsersHandler)
		r.PUT("/user/:id/role", middlewares.AuthMiddleware(erpContext, false), commonHandler.UpdateRoleHandler)
		r.DELETE("/user/:id", middlewares.AuthMiddleware(erpContext, false), commonHandler.DeleteUserHandler)
		r.GET("/members", middlewares.AuthMiddleware(erpContext, false), commonHandler.GetMembersHandler)
		r.GET("/roles", middlewares.AuthMiddleware(erpContext, false), commonHandler.GetRolesHandler)
		r.GET("/accept-invitation/:token", commonHandler.AcceptMemberInvitationHandler)
		r.POST("/invite-member", middlewares.AuthMiddleware(erpContext, false), middlewares.RbacUserMiddleware(erpContext, false, []string{"cooperative:cooperative_member:invite"}), commonHandler.InviteMemberHandler)
		r.POST("/invite-user", middlewares.AuthMiddleware(erpContext, false), middlewares.RbacUserMiddleware(erpContext, false, []string{"company:user:invite"}), commonHandler.InviteUserHandler)
		r.GET("/invited", middlewares.AuthMiddleware(erpContext, false), middlewares.RbacUserMiddleware(erpContext, false, []string{"cooperative:cooperative_member:invite"}), commonHandler.InvitedHandler)
		r.DELETE("/invited/:id", middlewares.AuthMiddleware(erpContext, false), middlewares.RbacUserMiddleware(erpContext, false, []string{"cooperative:cooperative_member:invite"}), commonHandler.DeleteInvitedHandler)
		r.GET("/setting", middlewares.AuthMiddleware(erpContext, false), commonHandler.CompanySettingHandler)
		r.PUT("/setting", middlewares.AuthMiddleware(erpContext, false), middlewares.RbacSuperAdminMiddleware(erpContext), commonHandler.UpdateCompanySettingHandler)
		r.PUT("/cooperative/setting", middlewares.AuthMiddleware(erpContext, false), middlewares.RbacSuperAdminMiddleware(erpContext), commonHandler.UpdateCooperativeSettingHandler)
	}
}
