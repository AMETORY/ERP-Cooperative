package handlers

import (
	"net/http"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/finance"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type TaxHandler struct {
	ctx        *context.ERPContext
	financeSrv *finance.FinanceService
}

func NewTaxHandler(ctx *context.ERPContext) *TaxHandler {
	financeSrv, ok := ctx.FinanceService.(*finance.FinanceService)
	if !ok {
		panic("finance service is not found")
	}
	return &TaxHandler{
		ctx:        ctx,
		financeSrv: financeSrv,
	}
}

func (t *TaxHandler) GetTaxHandler(c *gin.Context) {
	taxes, err := t.financeSrv.TaxService.GetTaxes(c.Request, c.Query("search"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": taxes, "message": "Taxes retrieved successfully"})
}

func (t *TaxHandler) GetTaxByIdHandler(c *gin.Context) {
	id := c.Param("id")
	tax, err := t.financeSrv.TaxService.GetTaxByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tax, "message": "Tax retrieved successfully"})
}

func (t *TaxHandler) CreateTaxHandler(c *gin.Context) {
	var input models.TaxModel
	err := c.ShouldBindBodyWithJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)
	userID := c.MustGet("userID").(string)
	input.CompanyID = &companyID
	input.UserID = &userID
	err = t.financeSrv.TaxService.CreateTax(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tax created successfully"})
}

func (t *TaxHandler) UpdateTaxHandler(c *gin.Context) {
	var input models.TaxModel
	err := c.ShouldBindBodyWithJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	_, err = t.financeSrv.TaxService.GetTaxByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	t.financeSrv.TaxService.UpdateTax(id, &input)
	c.JSON(http.StatusOK, gin.H{"message": "Tax updated successfully"})
}

func (t *TaxHandler) DeleteTaxHandler(c *gin.Context) {
	id := c.Param("id")
	err := t.financeSrv.TaxService.DeleteTax(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tax deleted successfully"})
}
