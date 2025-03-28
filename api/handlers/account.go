package handlers

import (
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/finance"
	"github.com/AMETORY/ametory-erp-modules/finance/account"
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
