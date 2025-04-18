package routes

import (
	"ametory-cooperative/api/handlers"
	"ametory-cooperative/api/middlewares"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupAssetRoutes(r *gin.RouterGroup, erpContext *context.ERPContext) {
	handler := handlers.NewAssetHandler(erpContext)
	assetGroup := r.Group("/asset")
	assetGroup.Use(middlewares.AuthMiddleware(erpContext, false))
	{
		assetGroup.GET("/list", handler.GetAssetListHandler)
		assetGroup.GET("/:id", handler.GetAssetHandler)
		assetGroup.GET("/:id/preview", handler.PreviewHandler)
		assetGroup.PUT("/:id/activate", handler.ActivateHandler)
		assetGroup.PUT("/:id/apply/:itemId", handler.DepreciationApplyHandler)
		assetGroup.POST("/create", handler.CreateAssetHandler)
		assetGroup.PUT("/:id", handler.UpdateAssetHandler)
		assetGroup.DELETE("/:id", handler.DeleteAssetHandler)
	}

}
