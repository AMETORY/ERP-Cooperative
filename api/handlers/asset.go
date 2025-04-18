package handlers

import (
	"ametory-cooperative/objects"
	"time"

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

	a.financeSrv.AssetService.GetDepreciation(asset)

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

func (a *AssetHandler) PreviewHandler(c *gin.Context) {
	id := c.Param("id")
	asset, err := a.financeSrv.AssetService.GetAssetByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	mode := c.Query("mode")
	isMonthly := c.DefaultQuery("is_monthly", "false")
	asset.DepreciationMethod = mode
	asset.IsMonthly = isMonthly == "true"
	preview, err := a.financeSrv.AssetService.PreviewCosts(asset)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": preview})
}

func (a *AssetHandler) ActivateHandler(c *gin.Context) {
	id := c.Param("id")

	input := objects.DepreciationRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	asset, err := a.financeSrv.AssetService.GetAssetByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("userID").(string)
	asset.DepreciationMethod = input.DepreciationMethod
	asset.AccountCurrentAssetID = &input.AccountCurrentAssetID
	asset.AccountFixedAssetID = &input.AccountFixedAssetID
	asset.AccountDepreciationID = &input.AccountDepreciationID
	asset.AccountAccumulatedDepreciationID = &input.AccountAccumulatedDepreciationID
	asset.IsMonthly = input.IsMonthly
	asset.Date = input.Date
	asset.Depreciations = input.DepreciationCosts
	if err := a.financeSrv.AssetService.ActivateAsset(asset, asset.Date, userID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": nil, "message": "Asset activated successfully"})
}

func (a *AssetHandler) DepreciationApplyHandler(c *gin.Context) {
	id := c.Param("id")
	itemId := c.Param("itemId")

	asset, err := a.financeSrv.AssetService.GetAssetByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("userID").(string)
	err = a.financeSrv.AssetService.DepreciationApply(asset, itemId, time.Now(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": nil, "message": "Asset depreciation applied successfully"})
}
