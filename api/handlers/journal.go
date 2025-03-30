package handlers

import (
	"ametory-cooperative/objects"
	"net/http"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/finance"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type JournalHandler struct {
	ctx        *context.ERPContext
	financeSrv *finance.FinanceService
}

func NewJournalHandler(ctx *context.ERPContext) *JournalHandler {
	financeSrv, ok := ctx.FinanceService.(*finance.FinanceService)
	if !ok {
		panic("finance service is not found")
	}
	return &JournalHandler{
		ctx:        ctx,
		financeSrv: financeSrv,
	}
}

func (h *JournalHandler) ListJournalsHandler(c *gin.Context) {
	journals, err := h.financeSrv.JournalService.GetJournals(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": journals})
}

func (h *JournalHandler) GetJournalHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	id := c.Param("id")
	journal, err := h.financeSrv.JournalService.GetJournal(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": journal})
}

func (h *JournalHandler) CreateJournalHandler(c *gin.Context) {
	var input models.JournalModel
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)
	userID := c.MustGet("userID").(string)
	input.CompanyID = &companyID
	input.UserID = &userID
	if err := h.financeSrv.JournalService.CreateJournal(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Journal created successfully", "id": input.ID})
}

func (h *JournalHandler) UpdateJournalHandler(c *gin.Context) {
	id := c.Param("id")
	var input models.JournalModel
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	journal, err := h.financeSrv.JournalService.GetJournal(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.financeSrv.JournalService.UpdateJournal(journal.ID, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Journal updated successfully"})
}

func (h *JournalHandler) DeleteJournalHandler(c *gin.Context) {
	id := c.Param("id")
	if err := h.financeSrv.JournalService.DeleteJournal(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Journal deleted successfully"})
}

func (h *JournalHandler) AddTransactionHandler(c *gin.Context) {
	id := c.Param("id")

	journal, err := h.financeSrv.JournalService.GetJournal(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var input objects.TransactionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)
	userID := c.MustGet("userID").(string)
	var transaction models.TransactionModel = models.TransactionModel{
		SourceID:      &input.SourceID,
		DestinationID: &input.DestinationID,
		Description:   input.Description,
		Date:          input.Date,
		CompanyID:     &companyID,
		UserID:        &userID,
		IsIncome:      input.IsIncome,
		IsExpense:     input.IsExpense,
		IsEquity:      input.IsEquity,
		IsTransfer:    input.IsTransfer,
		AccountID:     input.AccountID,
		Credit:        input.Credit,
		Debit:         input.Debit,
	}

	err = h.financeSrv.JournalService.AddTransaction(journal.ID, &transaction, input.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction added successfully"})
}
