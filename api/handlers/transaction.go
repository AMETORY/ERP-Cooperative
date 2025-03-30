package handlers

import (
	"ametory-cooperative/objects"
	"net/http"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/finance"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	ctx        *context.ERPContext
	financeSrv *finance.FinanceService
}

func NewTransactionHandler(ctx *context.ERPContext) *TransactionHandler {
	financeSrv, ok := ctx.FinanceService.(*finance.FinanceService)
	if !ok {
		panic("FinanceService is not instance of finance.FinanceService")
	}
	return &TransactionHandler{
		ctx:        ctx,
		financeSrv: financeSrv,
	}
}

func (h *TransactionHandler) ListTransactions(c *gin.Context) {
	transactions, err := h.financeSrv.TransactionService.GetTransactions(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": transactions})
}

func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	id := c.Param("id")
	transaction, err := h.financeSrv.TransactionService.GetTransactionById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transaction})
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
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
	if err := h.financeSrv.TransactionService.CreateTransaction(&transaction, input.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transaction created successfully"})
}

func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	id := c.Param("id")
	var input objects.TransactionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.financeSrv.TransactionService.GetTransactionById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	transaction.Amount = input.Amount
	transaction.Date = input.Date
	transaction.Description = input.Description
	if err := h.financeSrv.TransactionService.UpdateTransaction(id, transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transaction updated successfully"})
}

func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	id := c.Param("id")
	if err := h.financeSrv.TransactionService.DeleteTransaction(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}
