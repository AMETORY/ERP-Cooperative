package handlers

import (
	"strconv"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/shared/indonesia_regional"
	"github.com/AMETORY/ametory-erp-modules/thirdparty/google"
	"github.com/gin-gonic/gin"
)

type RegionalHandler struct {
	ctx             *context.ERPContext
	indonesiaRegSrv *indonesia_regional.IndonesiaRegService
}

func NewRegionalHandler(ctx *context.ERPContext) *RegionalHandler {
	return &RegionalHandler{
		ctx:             ctx,
		indonesiaRegSrv: ctx.IndonesiaRegService.(*indonesia_regional.IndonesiaRegService),
	}
}

func (h *RegionalHandler) GetProvinces(c *gin.Context) {
	search, _ := c.GetQuery("search")
	provinces, err := h.indonesiaRegSrv.GetProvinces(search)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"message": "Provinces retrieved successfully", "data": provinces})
}

func (h *RegionalHandler) GetRegencies(c *gin.Context) {
	var proviceID *string
	var search string
	provinceIDstring, ok := c.GetQuery("province_id")
	if ok {
		proviceID = &provinceIDstring
	}
	search, _ = c.GetQuery("search")
	regencies, err := h.indonesiaRegSrv.GetRegencies(proviceID, search)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"message": "Regencies retrieved successfully", "data": regencies})
}

func (h *RegionalHandler) GetDistricts(c *gin.Context) {
	var regencyID *string
	var search string
	regencyIDstring, ok := c.GetQuery("regency_id")
	if ok {
		regencyID = &regencyIDstring
	}
	search, _ = c.GetQuery("search")
	districts, err := h.indonesiaRegSrv.GetDistricts(regencyID, search)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"message": "Districts retrieved successfully", "data": districts})
}

func (h *RegionalHandler) GetVillages(c *gin.Context) {
	var districtID *string
	var search string
	districtIDstring, ok := c.GetQuery("district_id")
	if ok {
		districtID = &districtIDstring
	}
	search, _ = c.GetQuery("search")
	villages, err := h.indonesiaRegSrv.GetVillages(districtID, search)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"message": "Villages retrieved successfully", "data": villages})
}

func (h *RegionalHandler) SearchGooglePlaceHandler(c *gin.Context) {
	search := c.DefaultQuery("search", "")
	if search == "" {
		c.JSON(400, gin.H{"error": "search query is required"})
	}

	googleSrv, ok := h.ctx.ThirdPartyServices["google"]
	if !ok {
		c.JSON(500, gin.H{"error": "google service not found"})
	}
	resp, err := googleSrv.(*google.GoogleAPIService).SearchPlace(search)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}

	c.JSON(200, gin.H{"message": "Search Google Place successfully", "data": resp})
}
func (h *RegionalHandler) SearchGooglePlaceByCoordinateHandler(c *gin.Context) {
	latitude := c.DefaultQuery("latitude", "")
	if latitude == "" {
		c.JSON(400, gin.H{"error": "latitude query is required"})
	}
	longitude := c.DefaultQuery("longitude", "")
	if longitude == "" {
		c.JSON(400, gin.H{"error": "longitude query is required"})
	}

	googleSrv, ok := h.ctx.ThirdPartyServices["google"]
	if !ok {
		c.JSON(500, gin.H{"error": "google service not found"})
	}
	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	lon, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resp, err := googleSrv.(*google.GoogleAPIService).SearchPlaceByCoordinate(lat, lon, 5, 50)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}

	c.JSON(200, gin.H{"message": "Search Google Place successfully", "data": resp})
}
