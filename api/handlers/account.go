package handlers

import (
	"net/http"
	"time"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/finance"
	"github.com/AMETORY/ametory-erp-modules/finance/account"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	ctx        *context.ERPContext
	financeSrv *finance.FinanceService
}

func NewAccountHandler(ctx *context.ERPContext) *AccountHandler {
	financeSrv, ok := ctx.FinanceService.(*finance.FinanceService)
	if !ok {
		panic("FinanceService is not instance of finance.FinanceService")
	}
	return &AccountHandler{
		ctx:        ctx,
		financeSrv: financeSrv,
	}
}

func (a *AccountHandler) GetAccountTypesHandler(c *gin.Context) {
	types := a.financeSrv.AccountService.GetTypes()
	c.JSON(200, gin.H{"data": types})
}
func (a *AccountHandler) GetCodeHandler(c *gin.Context) {
	typeAccount := c.Query("type")
	var account models.AccountModel
	a.ctx.DB.Order("created_at desc").First(&account, "type = ?", typeAccount)

	c.JSON(200, gin.H{"last_code": account.Code})
}
func (a *AccountHandler) GetChartOfAccounts(c *gin.Context) {
	coa := account.GenericChartOfAccount
	template := c.Query("template")
	if template == "cooperative" {
		coa = account.CooperationAccountsTemplate
	}
	if template == "islamic-cooperative" {
		coa = account.IslamicCooperationAccountsTemplate
	}

	c.JSON(200, gin.H{"data": coa})
}

func (h *AccountHandler) CreateAccountHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to create an account
	if h.financeSrv == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "inventory service is not initialized"})
	}
	var data models.AccountModel
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)
	userID := c.MustGet("userID").(string)
	data.CompanyID = &companyID
	data.UserID = &userID
	data.IsDeletable = true
	err = h.financeSrv.AccountService.CreateAccount(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account created successfully"})
}

func (h *AccountHandler) GetAccountHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to get an account
	if h.financeSrv == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "distributor service is not initialized"})
	}
	search, _ := c.GetQuery("search")
	data, err := h.financeSrv.AccountService.GetAccounts(*c.Request, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newItems := make([]models.AccountModel, 0)
	items := data.Items.(*[]models.AccountModel)
	now := time.Now().AddDate(0, 0, 1)
	for _, item := range *items {
		debit, credit, _ := h.financeSrv.ReportService.GetAccountBalance(item.ID, nil, &now)
		switch item.Type {
		case models.EXPENSE, models.COST, models.CONTRA_LIABILITY, models.CONTRA_EQUITY, models.CONTRA_REVENUE:
			item.Balance = debit - credit
		case models.LIABILITY, models.EQUITY, models.REVENUE, models.INCOME, models.CONTRA_ASSET, models.CONTRA_EXPENSE:
			item.Balance = credit - debit
		case models.ASSET:
			item.Balance = debit - credit
		}
		newItems = append(newItems, item)
	}

	*items = newItems

	c.JSON(http.StatusOK, gin.H{"message": "Account retrieved successfully", "data": data})
}

func (h *AccountHandler) GetAccountReportHandler(c *gin.Context) {
	id := c.Param("id")
	data, err := h.financeSrv.ReportService.GenerateAccountReport(id, *c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account report retrieved successfully", "data": data})
}
func (h *AccountHandler) GetAccountByIdHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to get an account by ID
	if h.financeSrv == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "distributor service is not initialized"})
	}
	id := c.Param("id")
	data, err := h.financeSrv.AccountService.GetAccountByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account retrieved successfully", "data": data})
}

func (h *AccountHandler) UpdateAccountHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to update an account
	var data models.AccountModel
	err := c.ShouldBindBodyWithJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	_, err = h.financeSrv.AccountService.GetAccountByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	h.financeSrv.AccountService.UpdateAccount(id, &data)
	c.JSON(http.StatusOK, gin.H{"message": "Account updated successfully"})
}

func (h *AccountHandler) DeleteAccountHandler(c *gin.Context) {
	h.ctx.Request = c.Request
	// Implement logic to delete an account
	if h.financeSrv == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "distributor service is not initialized"})
	}
	id := c.Param("id")
	err := h.financeSrv.AccountService.DeleteAccount(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
