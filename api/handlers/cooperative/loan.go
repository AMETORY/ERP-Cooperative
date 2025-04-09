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

	userID := c.MustGet("userID").(string)
	companyID := c.MustGet("companyID").(string)
	input.CompanyID = &companyID
	input.UserID = &userID

	err = h.coopertiveSrv.LoanApplicationService.CreateLoan(&input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Loan created successfully", "data": input})
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

func (h *LoanApplicationHandler) ApprovalHandler(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		ApprovalStatus string `json:"approval_status" binding:"required"`
		Remarks        string `json:"remarks" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(string)

	err := h.coopertiveSrv.LoanApplicationService.ApprovalLoan(id, userID, input.ApprovalStatus, input.Remarks)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Loan approved successfully"})
}

func (h *LoanApplicationHandler) DisbursementHandler(c *gin.Context) {
	var input struct {
		AccountAssetID string `json:"account_asset_id" binding:"required"`
		Remarks        string `json:"remarks" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return

	}
	accountAssetID := input.AccountAssetID
	id := c.Param("id")
	loan, err := h.coopertiveSrv.LoanApplicationService.GetLoanByID(id, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(models.UserModel)
	err = h.coopertiveSrv.LoanApplicationService.DisburseLoan(loan, &accountAssetID, &user, input.Remarks)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Loan disbursed successfully"})
}

func (h *LoanApplicationHandler) PaymentHandler(c *gin.Context) {
	id := c.Param("id")
	loan, err := h.coopertiveSrv.LoanApplicationService.GetLoanByID(id, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var input models.InstallmentPayment
	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	input.MemberID = loan.MemberID
	userID := c.MustGet("userID").(string)
	err = h.coopertiveSrv.LoanApplicationService.CreatePayment(&input, loan, &userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Loan payment added successfully", "data": input})
}
