package handlers

import (
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/finance"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	ctx        *context.ERPContext
	financeSrv *finance.FinanceService
}

func NewAssetHandler(ctx *context.ERPContext) *AssetHandler {
	financeSrv, ok := ctx.FinanceService.(*finance.FinanceService)
	if !ok {
		panic("FinanceService is not instance of finance.FinanceService")
	}
	return &AssetHandler{
		ctx:        ctx,
		financeSrv: financeSrv,
	}
}

func (a *AssetHandler) GetAssetListHandler(c *gin.Context) {
	assets, err := a.financeSrv.AssetService.GetAssets(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": assets})
}

func (a *AssetHandler) GetAssetHandler(c *gin.Context) {
	id := c.Param("id")
	asset, err := a.financeSrv.AssetService.GetAssetByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": asset})
}

func (a *AssetHandler) CreateAssetHandler(c *gin.Context) {
	var asset models.AssetModel
	if err := c.BindJSON(&asset); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("userID").(string)
	companyID := c.MustGet("companyID").(string)
	asset.UserID = &userID
	asset.CompanyID = &companyID
	if err := a.financeSrv.AssetService.CreateAsset(&asset); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": asset})
}

func (a *AssetHandler) UpdateAssetHandler(c *gin.Context) {
	id := c.Param("id")
	var asset models.AssetModel
	if err := c.BindJSON(&asset); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := a.financeSrv.AssetService.UpdateAsset(id, &asset); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": asset})
}

func (a *AssetHandler) DeleteAssetHandler(c *gin.Context) {
	id := c.Param("id")
	if err := a.financeSrv.AssetService.DeleteAsset(id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": nil})
}
