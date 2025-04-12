package handlers

import (
	"ametory-cooperative/app_models"
	"ametory-cooperative/objects"
	"net/http"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/finance"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/AMETORY/ametory-erp-modules/utils"
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

func (r *ReportHandler) GenerateClosingBookHandler(c *gin.Context) {
	input := struct {
		Notes           string  `json:"notes" binding:"required"`
		RetainEarningId string  `json:"retain_earning_id" binding:"required"`
		TaxPercentage   float64 `json:"tax_percentage"`
		TaxPayableId    *string `json:"tax_payable_id"`
		TaxExpenseID    *string `json:"tax_expense_id"`
	}{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Params.ByName("id")
	report, err := r.financeSrv.ReportService.GetClosingBookByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var setting app_models.CustomSettingModel
	err = r.ctx.DB.Where("id = ?", c.GetHeader("ID-Company")).First(&setting).Error
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(string)

	err = r.financeSrv.ReportService.GenerateClosingBook(
		report,
		setting.CashflowGroupSetting,
		userID,
		input.Notes,
		input.RetainEarningId,
		input.TaxPayableId,
		input.TaxExpenseID,
		input.TaxPercentage,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "closing book generated successfully"})
}

func (r *ReportHandler) GetClosingBookDetailHandler(c *gin.Context) {
	id := c.Params.ByName("id")
	report, err := r.financeSrv.ReportService.GetClosingBookByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var setting app_models.CustomSettingModel
	err = r.ctx.DB.Where("id = ?", c.GetHeader("ID-Company")).First(&setting).Error
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	companyID := c.MustGet("companyID").(string)

	cashflowReport := models.CashFlowReport{}
	cashflowReport.StartDate = report.StartDate
	cashflowReport.EndDate = report.EndDate
	cashflowReport.CompanyID = companyID
	cashflowReport.Operating = setting.CashflowGroupSetting.Operating
	cashflowReport.Investing = setting.CashflowGroupSetting.Investing
	cashflowReport.Financing = setting.CashflowGroupSetting.Financing

	if report.CashFlow == nil {
		cashFlow, err := r.financeSrv.ReportService.GenerateCashFlowReport(cashflowReport)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		report.CashFlow = cashFlow
	}
	if report.TrialBalance == nil {
		trialBalance, err := r.financeSrv.ReportService.TrialBalanceReport(models.GeneralReport{
			StartDate: report.StartDate,
			EndDate:   report.EndDate,
			CompanyID: companyID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		report.TrialBalance = trialBalance

	}

	if report.ProfitLoss == nil {
		profitLoss, err := r.financeSrv.ReportService.GenerateProfitLossReport(models.GeneralReport{
			StartDate: report.StartDate,
			EndDate:   report.EndDate,
			CompanyID: companyID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		report.ProfitLoss = profitLoss
	}

	if report.BalanceSheet == nil {
		balanceSheet, err := r.financeSrv.ReportService.GenerateBalanceSheet(models.GeneralReport{
			StartDate: report.StartDate,
			EndDate:   report.EndDate,
			CompanyID: companyID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		report.BalanceSheet = balanceSheet
	}

	if report.CapitalChange == nil {
		capitalChange, err := r.financeSrv.ReportService.GenerateCapitalChangeReport(models.GeneralReport{
			StartDate: report.StartDate,
			EndDate:   report.EndDate,
			CompanyID: companyID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		report.CapitalChange = capitalChange
	}

	c.JSON(http.StatusOK, gin.H{"message": "closing book retrieved successfully", "data": report})
}
func (r *ReportHandler) GetClosingBooksHandler(c *gin.Context) {

	report, err := r.financeSrv.ReportService.GetClosingBook(*c.Request, c.Query("search"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	c.JSON(http.StatusOK, gin.H{"message": "closing book retrieved successfully", "data": report})

}

func (r *ReportHandler) CreateClosingBookHandler(c *gin.Context) {
	input := models.ClosingBook{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyID := c.MustGet("companyID").(string)
	userID := c.MustGet("userID").(string)
	input.CompanyID = &companyID
	input.UserID = &userID
	input.Status = "DRAFT"
	input.ID = utils.Uuid()

	err := r.financeSrv.ReportService.CreateClosingBook(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Closing book created successfully", "data": input})
}

func (r *ReportHandler) TrialBalanceHandler(c *gin.Context) {
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

	report, err := r.financeSrv.ReportService.TrialBalanceReport(models.GeneralReport{
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		CompanyID: companyID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "trial balance retrieved successfully", "data": report})
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
