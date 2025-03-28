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
		r.GET("/members", middlewares.AuthMiddleware(erpContext, false), commonHandler.GetMembersHandler)
		r.GET("/roles", middlewares.AuthMiddleware(erpContext, false), commonHandler.GetRolesHandler)
		r.GET("/accept-invitation/:token", commonHandler.AcceptMemberInvitationHandler)
		r.POST("/invite-member", middlewares.AuthMiddleware(erpContext, false), middlewares.RbacUserMiddleware(erpContext, false, []string{"cooperative:cooperative_member:create"}), commonHandler.InviteMemberHandler)
		r.GET("/invited", middlewares.AuthMiddleware(erpContext, false), middlewares.RbacUserMiddleware(erpContext, false, []string{"cooperative:cooperative_member:create"}), commonHandler.InvitedHandler)
		r.DELETE("/invited/:id", middlewares.AuthMiddleware(erpContext, false), middlewares.RbacUserMiddleware(erpContext, false, []string{"cooperative:cooperative_member:create"}), commonHandler.DeleteInvitedHandler)
		r.GET("/setting", middlewares.AuthMiddleware(erpContext, false), middlewares.RbacSuperAdminMiddleware(erpContext), commonHandler.CompanySettingHandler)
		r.PUT("/setting", middlewares.AuthMiddleware(erpContext, false), middlewares.RbacSuperAdminMiddleware(erpContext), commonHandler.UpdateCompanySettingHandler)
	}
}
