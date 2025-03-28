package routes

import (
	"ametory-cooperative/api/handlers"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/gin-gonic/gin"
)

func SetupRegionalRoutes(r *gin.RouterGroup, ctx *context.ERPContext) {
	regionalHandler := handlers.NewRegionalHandler(ctx)

	regionalGroup := r.Group("/regional")
	regionalGroup.GET("/province", regionalHandler.GetProvinces)
	regionalGroup.GET("/regency", regionalHandler.GetRegencies)
	regionalGroup.GET("/district", regionalHandler.GetDistricts)
	regionalGroup.GET("/village", regionalHandler.GetVillages)
	regionalGroup.GET("/place", regionalHandler.SearchGooglePlaceHandler)
	regionalGroup.GET("/place-by-coordinate", regionalHandler.SearchGooglePlaceByCoordinateHandler)
}
