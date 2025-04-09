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
	loanHandler := cooperative_handler.NewLoanApplicationHandler(ctx)
	loanGroup := r.Group("/cooperative/loan")
	loanGroup.Use(middlewares.AuthMiddleware(ctx, false))
	{
		loanGroup.GET("/list", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:loan_application:read"}), loanHandler.GetLoansHandler)
		loanGroup.GET("/:id", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:loan_application:read"}), loanHandler.GetLoanHandler)
		loanGroup.POST("/create", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:loan_application:create"}), loanHandler.CreateLoanHandler)
		loanGroup.PUT("/:id", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:loan_application:update"}), loanHandler.UpdateLoanHandler)
		loanGroup.PUT("/:id/approval", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:loan_application:update"}), loanHandler.ApprovalHandler)
		loanGroup.PUT("/:id/disbursement", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:loan_application:update"}), loanHandler.DisbursementHandler)
		loanGroup.PUT("/:id/payment", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:loan_application:update"}), loanHandler.PaymentHandler)
		loanGroup.DELETE("/:id", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:loan_application:delete"}), loanHandler.DeleteLoanHandler)
	}

	savingHandler := cooperative_handler.NewSavingHandler(ctx)
	savingGroup := r.Group("/cooperative/saving")
	savingGroup.Use(middlewares.AuthMiddleware(ctx, false))
	{
		savingGroup.GET("/list", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:saving:read"}), savingHandler.GetSavingsHandler)
		savingGroup.GET("/:id", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:saving:read"}), savingHandler.GetSavingHandler)
		savingGroup.POST("/create", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:saving:create"}), savingHandler.CreateSavingHandler)
		savingGroup.PUT("/:id", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:saving:update"}), savingHandler.UpdateSavingHandler)
		savingGroup.DELETE("/:id", middlewares.RbacUserMiddleware(ctx, false, []string{"cooperative:saving:delete"}), savingHandler.DeleteSavingHandler)
	}

}
