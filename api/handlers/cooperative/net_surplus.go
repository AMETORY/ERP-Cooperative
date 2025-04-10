package cooperative_handler

import (
	"time"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/cooperative"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type NetSurplusHandler struct {
	ctx           *context.ERPContext
	coopertiveSrv *cooperative.CooperativeService
}

func NewNetSurplusHandler(ctx *context.ERPContext) *NetSurplusHandler {
	cooperativeSrv, ok := ctx.CooperativeService.(*cooperative.CooperativeService)
	if !ok {
		panic("CooperativeService is not found")
	}
	return &NetSurplusHandler{
		ctx:           ctx,
		coopertiveSrv: cooperativeSrv,
	}
}

func (h *NetSurplusHandler) GetNetSurplusListHandler(c *gin.Context) {
	netSurplus, err := h.coopertiveSrv.NetSurplusService.GetNetSurplusList(*c.Request, c.Query("search"), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": netSurplus})
}

func (h *NetSurplusHandler) GetNetSurplusHandler(c *gin.Context) {
	id := c.Param("id")
	netSurplus, err := h.coopertiveSrv.NetSurplusService.GetNetSurplusByID(id, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": netSurplus})
}

func (h *NetSurplusHandler) CreateNetSurplusHandler(c *gin.Context) {
	var input models.NetSurplusModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("userID").(string)
	companyID := c.MustGet("companyID").(string)
	input.CompanyID = &companyID
	input.UserID = &userID

	var company models.CompanyModel
	err = h.ctx.DB.Where("id = ?", companyID).First(&company).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = h.coopertiveSrv.NetSurplusService.CreateNetSurplus(&input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "NetSurplus created successfully", "data": input})
}

func (h *NetSurplusHandler) UpdateNetSurplusHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.NetSurplusModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = h.coopertiveSrv.NetSurplusService.UpdateNetSurplus(id, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "NetSurplus updated successfully"})
}
func (h *NetSurplusHandler) DistributeNetSurplusHandler(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		SourceID    string                        `json:"source_id" binding:"required"`
		Allocations []models.NetSurplusAllocation `json:"allocations" binding:"required"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	netSurplus, err := h.coopertiveSrv.NetSurplusService.GetNetSurplusByID(id, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("userID").(string)
	err = h.coopertiveSrv.NetSurplusService.Distribute(netSurplus, input.SourceID, input.Allocations, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "NetSurplus distributed successfully"})
}
func (h *NetSurplusHandler) DisbursementeNetSurplusHandler(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Date          time.Time                 `json:"date" binding:"required"`
		DestinationID string                    `json:"destination_id" binding:"required"`
		Members       []models.NetSurplusMember `json:"members" binding:"required"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	netSurplus, err := h.coopertiveSrv.NetSurplusService.GetNetSurplusByID(id, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("userID").(string)
	err = h.coopertiveSrv.NetSurplusService.Disbursement(input.Date, input.Members, netSurplus, input.DestinationID, userID, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "NetSurplus disbursement successfully"})
}

func (h *NetSurplusHandler) DeleteNetSurplusHandler(c *gin.Context) {
	id := c.Param("id")
	err := h.coopertiveSrv.NetSurplusService.DeleteNetSurplus(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "NetSurplus deleted successfully"})
}
