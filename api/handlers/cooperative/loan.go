package cooperative_handler

import (
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/cooperative"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type LoanApplicationHandler struct {
	ctx           *context.ERPContext
	coopertiveSrv *cooperative.CooperativeService
}

func NewLoanApplicationHandler(ctx *context.ERPContext) *LoanApplicationHandler {
	cooperativeSrv, ok := ctx.CooperativeService.(*cooperative.CooperativeService)
	if !ok {
		panic("CooperativeService is not found")
	}
	return &LoanApplicationHandler{
		ctx:           ctx,
		coopertiveSrv: cooperativeSrv,
	}
}

func (h *LoanApplicationHandler) GetLoansHandler(c *gin.Context) {
	loans, err := h.coopertiveSrv.LoanApplicationService.GetLoans(*c.Request, c.Query("search"), nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": loans})
}

func (h *LoanApplicationHandler) GetLoanHandler(c *gin.Context) {
	id := c.Param("id")
	loan, err := h.coopertiveSrv.LoanApplicationService.GetLoanByID(id, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": loan})
}

func (h *LoanApplicationHandler) CreateLoanHandler(c *gin.Context) {
	var input models.LoanApplicationModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = h.coopertiveSrv.LoanApplicationService.CreateLoan(&input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Loan created successfully"})
}

func (h *LoanApplicationHandler) UpdateLoanHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.LoanApplicationModel
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = h.coopertiveSrv.LoanApplicationService.UpdateLoan(id, &input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Loan updated successfully"})
}

func (h *LoanApplicationHandler) DeleteLoanHandler(c *gin.Context) {
	id := c.Param("id")
	err := h.coopertiveSrv.LoanApplicationService.DeleteLoan(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Loan deleted successfully"})
}
