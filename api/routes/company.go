package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupCompanyRoutes(r *gin.RouterGroup, ctx *context.ERPContext) {
	companyHandler := handlers.NewCompanyHandler(ctx)
	companyGroup := r.Group("/company")
	companyGroup.Use(middlewares.AuthMiddleware(ctx, false))
	{
		companyGroup.GET("/categories", companyHandler.GetCategories)
		companyGroup.GET("/sectors", companyHandler.GetSectors)
		companyGroup.POST("/create", companyHandler.CreateCompany)
	}
}
