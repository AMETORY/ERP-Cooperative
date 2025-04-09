package handlers

import (
	"ametory-cooperative/app_models"
	"ametory-cooperative/objects"
	"net/http"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/finance"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	ctx        *context.ERPContext
	financeSrv *finance.FinanceService
}

func NewReportHandler(ctx *context.ERPContext) *ReportHandler {
	financeSrv, ok := ctx.FinanceService.(*finance.FinanceService)
	if !ok {
		panic("finance service is not found")
	}
	return &ReportHandler{
		ctx:        ctx,
		financeSrv: financeSrv,
	}
}

func (r *ReportHandler) GetCogsHandler(c *gin.Context) {
	input := objects.ReportRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)

	report, err := r.financeSrv.ReportService.GenerateCogsReport(models.GeneralReport{
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		CompanyID: companyID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "profit loss report retrieved successfully", "data": report})
}
func (r *ReportHandler) GetProfitLossHandler(c *gin.Context) {
	input := objects.ReportRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)

	report, err := r.financeSrv.ReportService.GenerateProfitLossReport(models.GeneralReport{
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		CompanyID: companyID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cogs report retrieved successfully", "data": report})
}
func (r *ReportHandler) GetBalanceSheetHandler(c *gin.Context) {
	input := objects.ReportRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)

	report, err := r.financeSrv.ReportService.GenerateBalanceSheet(models.GeneralReport{
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		CompanyID: companyID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "balance sheet retrieved successfully", "data": report})
}
func (r *ReportHandler) CapitalChangeHandler(c *gin.Context) {
	input := objects.ReportRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)

	report, err := r.financeSrv.ReportService.GenerateCapitalChangeReport(models.GeneralReport{
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		CompanyID: companyID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "capital change retrieved successfully", "data": report})
}

func (r *ReportHandler) CashFlowHandler(c *gin.Context) {
	input := objects.ReportRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var setting app_models.CustomSettingModel
	err := r.ctx.DB.Where("id = ?", c.GetHeader("ID-Company")).First(&setting).Error
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	companyID := c.MustGet("companyID").(string)

	cashflowReport := models.CashFlowReport{}
	cashflowReport.StartDate = input.StartDate
	cashflowReport.EndDate = input.EndDate
	cashflowReport.CompanyID = companyID
	cashflowReport.Operating = setting.CashflowGroupSetting.Operating
	cashflowReport.Investing = setting.CashflowGroupSetting.Investing
	cashflowReport.Financing = setting.CashflowGroupSetting.Financing

	report, err := r.financeSrv.ReportService.GenerateCashFlowReport(cashflowReport)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cash flow retrieved successfully", "data": report})
}
