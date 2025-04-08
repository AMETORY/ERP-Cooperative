package cooperative_handler

import (
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/cooperative"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type SavingHandler struct {
	ctx           *context.ERPContext
	coopertiveSrv *cooperative.CooperativeService
}

func NewSavingHandler(ctx *context.ERPContext) *SavingHandler {
	cooperativeSrv, ok := ctx.CooperativeService.(*cooperative.CooperativeService)
	if !ok {
		panic("CooperativeService is not found")
	}
	return &SavingHandler{
		ctx:           ctx,
		coopertiveSrv: cooperativeSrv,
	}
}

func (h *SavingHandler) GetSavingsHandler(c *gin.Context) {
	loans, err := h.coopertiveSrv.SavingService.GetSavings(*c.Request, c.Query("search"), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": loans})
}

func (h *SavingHandler) GetSavingHandler(c *gin.Context) {
	id := c.Param("id")
	loan, err := h.coopertiveSrv.SavingService.GetSavingByID(id, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": loan})
}

func (h *SavingHandler) CreateSavingHandler(c *gin.Context) {
	var input models.SavingModel
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
	h.ctx.DB.Where("id = ?", companyID).First(&company)
	input.Company = &company

	err = h.coopertiveSrv.SavingService.CreateSaving(&input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = h.coopertiveSrv.SavingService.CreateTransaction(input, false)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		err := h.coopertiveSrv.SavingService.DeleteSaving(input.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		return
	}
	c.JSON(201, gin.H{"message": "Saving created successfully"})
}

func (h *SavingHandler) UpdateSavingHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.SavingModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = h.coopertiveSrv.SavingService.UpdateSaving(id, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Saving updated successfully"})
}

func (h *SavingHandler) DeleteSavingHandler(c *gin.Context) {
	id := c.Param("id")
	err := h.coopertiveSrv.SavingService.DeleteSaving(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Saving deleted successfully"})
}
