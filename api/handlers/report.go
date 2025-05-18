package handlers

import (
	"ametory-cooperative/app_models"
	"ametory-cooperative/objects"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/AMETORY/ametory-erp-modules/contact"
	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/finance"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/AMETORY/ametory-erp-modules/utils"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type ReportHandler struct {
	ctx        *context.ERPContext
	financeSrv *finance.FinanceService
	contactSev *contact.ContactService
}

var borderAll = []excelize.Border{
	{
		Type:  "left",
		Color: "#555555",
		Style: 1,
	},
	{
		Type:  "right",
		Color: "#555555",
		Style: 1,
	},
	{
		Type:  "top",
		Color: "#555555",
		Style: 1,
	},
	{
		Type:  "bottom",
		Color: "#555555",
		Style: 1,
	},
}
var borderHorizontal = []excelize.Border{
	{
		Type:  "left",
		Color: "#555555",
		Style: 1,
	},
	{
		Type:  "right",
		Color: "#555555",
		Style: 1,
	},
}

func NewReportHandler(ctx *context.ERPContext) *ReportHandler {
	financeSrv, ok := ctx.FinanceService.(*finance.FinanceService)
	if !ok {
		panic("finance service is not found")
	}
	contactSev, ok := ctx.ContactService.(*contact.ContactService)
	if !ok {
		panic("contact service is not found")
	}
	return &ReportHandler{
		ctx:        ctx,
		financeSrv: financeSrv,
		contactSev: contactSev,
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
		Notes           string  `json:"notes"`
		RetainEarningId string  `json:"retain_earning_id" binding:"required"`
		ProfitSummaryID string  `json:"profit_summary_id" binding:"required"`
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
		input.ProfitSummaryID,
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

func (r *ReportHandler) DeleteClosingBooklHandler(c *gin.Context) {
	id := c.Params.ByName("id")
	err := r.financeSrv.ReportService.DeleteClosingBook(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "closing book deleted successfully"})
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

	if report.ClosingSummary == nil {
		var summary models.ClosingSummary

		summary.TotalIncome = report.ProfitLoss.GrossProfit
		summary.TotalExpense = report.ProfitLoss.TotalExpense
		summary.NetIncome = report.ProfitLoss.NetProfit
		report.ClosingSummary = &summary
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

func (r *ReportHandler) GetProductSalesCustomersHandler(c *gin.Context) {
	input := objects.ReportRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)

	report, err := r.financeSrv.ReportService.GetProductSalesCustomers(companyID, input.StartDate, input.EndDate, input.ProductIDs, input.CustomerIDs)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	matrix := map[string]map[string][]models.ProductSalesCustomer{}
	productSet := map[string]models.ProductSalesCustomer{}
	contactSet := map[string]models.ProductSalesCustomer{}
	grandTotalQuantity := map[string]float64{}
	grandTotalAmount := map[string]float64{}

	for _, d := range report {
		if _, ok := matrix[d.ProductID]; !ok {
			matrix[d.ProductID] = map[string][]models.ProductSalesCustomer{}
		}
		matrix[d.ProductID][d.ContactID] = append(matrix[d.ProductID][d.ContactID], d)
		grandTotalQuantity[d.ProductID] += d.TotalQuantity
		grandTotalAmount[d.ProductID] += d.TotalPrice
		productSet[d.ProductID] = models.ProductSalesCustomer{
			ProductID:   d.ProductID,
			ProductCode: d.ProductCode,
			ProductName: d.ProductName,
			UnitCode:    d.UnitCode,
			UnitName:    d.UnitName,
		}
		contactSet[d.ContactID] = models.ProductSalesCustomer{
			ContactID:   d.ContactID,
			ContactCode: d.ContactCode,
			ContactName: d.ContactName,
		}
	}

	// sort productSet by product_code
	sortedProducts := make([]models.ProductSalesCustomer, 0, len(productSet))
	for _, v := range productSet {
		sortedProducts = append(sortedProducts, v)
	}
	sort.Slice(sortedProducts, func(i, j int) bool {
		return sortedProducts[i].ProductCode < sortedProducts[j].ProductCode
	})
	productSet = map[string]models.ProductSalesCustomer{}
	for _, v := range sortedProducts {
		productSet[v.ProductID] = v
	}

	// sort contactSet by contact_code
	sortedContacts := make([]models.ProductSalesCustomer, 0, len(contactSet))
	for _, v := range contactSet {
		sortedContacts = append(sortedContacts, v)
	}
	sort.Slice(sortedContacts, func(i, j int) bool {
		return sortedContacts[i].ContactCode < sortedContacts[j].ContactCode
	})
	contactSet = map[string]models.ProductSalesCustomer{}
	for _, v := range sortedContacts {
		contactSet[v.ContactID] = v
	}

	fmtCode := `_(#,##0_);_(\(#,##0\);_("-"??_);_(@_)`

	if input.IsDownload {
		file := excelize.NewFile()
		sheet1 := file.GetSheetName(0)

		headerStyle, err := file.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold: true,
				Size: 12,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#DCE6F1"}, // Soft blue
				Pattern: 1,
			},
		})
		boldStyle, _ := file.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold: true,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
		})
		boldStyleFormat, _ := file.NewStyle(&excelize.Style{
			CustomNumFmt: &fmtCode,
			Font: &excelize.Font{
				Bold: true,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
		})
		centerStyle, _ := file.NewStyle(&excelize.Style{
			CustomNumFmt: &fmtCode,
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
		})
		row := 1
		headers := []string{"Kode"}
		colWidth := []float64{15}
		for _, v := range productSet {

			headers = append(headers, v.ProductName)
			colWidth = append(colWidth, 20)
		}

		for i, header := range headers {
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", utils.NumToAlphabet(i+1), row), header)
			file.SetColWidth(sheet1, utils.NumToAlphabet(i+1), utils.NumToAlphabet(i+1), colWidth[i])
			// Apply styles: bold font, bigger font, center align, and soft blue background

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			file.SetCellStyle(sheet1, fmt.Sprintf("A%d", row), fmt.Sprintf("%s%d", utils.NumToAlphabet(len(headers)), row), headerStyle)

		}

		row++
		headers = []string{""}
		colWidth = []float64{15}
		for _, v := range productSet {
			if input.View == "quantity" && v.UnitCode != "" {
				headers = append(headers, fmt.Sprintf(`%s
(%s)`, v.ProductCode, v.UnitCode))
			} else {
				headers = append(headers, v.ProductCode)

			}
			colWidth = append(colWidth, 20)
		}

		for i, header := range headers {
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", utils.NumToAlphabet(i+1), row), header)
			file.SetColWidth(sheet1, utils.NumToAlphabet(i+1), utils.NumToAlphabet(i+1), colWidth[i])
			// Apply styles: bold font, bigger font, center align, and soft blue background

			file.SetCellStyle(sheet1, fmt.Sprintf("A%d", row), fmt.Sprintf("%s%d", utils.NumToAlphabet(len(headers)), row), headerStyle)

		}

		file.MergeCell(sheet1, "A1", "A2")
		row++

		for keyCustomer, v := range contactSet {
			file.SetCellValue(sheet1, fmt.Sprintf("A%d", row), v.ContactCode)
			file.SetCellStyle(sheet1, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), boldStyle)
			i := 1
			for key := range productSet {
				var value any
				if input.View == "quantity" {
					val, ok := matrix[key][keyCustomer]
					if ok && len(val) > 0 {
						value = val[0].TotalQuantity
					}
				}
				if input.View == "amount" {
					val, ok := matrix[key][keyCustomer]
					if ok && len(val) > 0 {
						value = val[0].TotalPrice
					}
				}

				file.SetCellValue(sheet1, fmt.Sprintf("%s%d", utils.NumToAlphabet(i+1), row), value)
				file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", utils.NumToAlphabet(i+1), row), fmt.Sprintf("%s%d", utils.NumToAlphabet(i+1), row), centerStyle)
				i++
			}
			row++
		}

		footers := []any{"Grand Total"}
		colWidth = []float64{15}
		for _, v := range productSet {
			if input.View == "quantity" {
				footers = append(footers, grandTotalQuantity[v.ProductID])
			} else {
				footers = append(footers, grandTotalAmount[v.ProductID])
			}
			colWidth = append(colWidth, 20)
		}

		for i, footer := range footers {
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", utils.NumToAlphabet(i+1), row), footer)
			file.SetColWidth(sheet1, utils.NumToAlphabet(i+1), utils.NumToAlphabet(i+1), colWidth[i])
			if i > 0 {
				file.SetCellStyle(sheet1, fmt.Sprintf("A%d", row), fmt.Sprintf("%s%d", utils.NumToAlphabet(len(footers)), row), boldStyleFormat)
			} else {
				file.SetCellStyle(sheet1, fmt.Sprintf("A%d", row), fmt.Sprintf("%s%d", utils.NumToAlphabet(len(footers)), row), boldStyle)
			}
		}

		var buf bytes.Buffer
		if err := file.Write(&buf); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write XLSX file"})
			return
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename=report.xlsx")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product sales customers report retrieved successfully", "data": matrix, "products": productSet, "contacts": contactSet, "grand_total_amount": grandTotalAmount, "grand_total_quantity": grandTotalQuantity})

}

func (r *ReportHandler) GetAccountReceivableLedgerHandler(c *gin.Context) {
	var input struct {
		Title      string    `json:"title,omitempty" form:"title"`
		StartDate  time.Time `json:"start_date,omitempty" form:"start_date"`
		EndDate    time.Time `json:"end_date,omitempty" form:"end_date"`
		ContactID  string    `json:"contact_id,omitempty" example:"contact_id"`
		IsDownload bool      `json:"is_download,omitempty" form:"is_download"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)
	// TODO: add date filter
	r.financeSrv.ReportService.SetContactService(r.contactSev)
	data, err := r.financeSrv.ReportService.GetAccountReceivableLedger(companyID, input.ContactID, input.StartDate, input.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if input.IsDownload {
		fmtCode := `_(#,##0_);_(\(#,##0\);_("-"??_);_(@_)`
		file := excelize.NewFile()
		sheet1 := file.GetSheetName(0)

		titleStyle, _ := file.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold: true,
				Size: 20,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
				Indent:     10,
			},
			Border: borderAll,
		})
		headerStyle, _ := file.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold: true,
				Size: 12,
			},

			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
				Indent:     1,
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#DCE6F1"}, // Soft blue
				Pattern: 1,
			},
			Border: borderAll,
		})
		heroStyle, _ := file.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold: true,
				Size: 12,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "left",
				Vertical:   "center",
				Indent:     1,
			},
			Border: borderAll,
		})
		heroStyleNormal, _ := file.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Size: 12,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "left",
				Vertical:   "center",
				Indent:     1,
			},
			Border: borderAll,
		})
		heroStyleWithFormat, _ := file.NewStyle(&excelize.Style{
			CustomNumFmt: &fmtCode,
			Font: &excelize.Font{
				Size: 12,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "left",
				Vertical:   "center",
				Indent:     1,
			},
			Border: borderAll,
		})
		boldStyle, _ := file.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold: true,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
				Indent:     1,
			},
			Border: borderAll,
		})
		boldStyleFormat, _ := file.NewStyle(&excelize.Style{
			CustomNumFmt: &fmtCode,
			Font: &excelize.Font{
				Bold: true,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
				Indent:     1,
			},
			Border: borderAll,
		})
		boldLeftStyleFormat, _ := file.NewStyle(&excelize.Style{
			CustomNumFmt: &fmtCode,
			Font: &excelize.Font{
				Bold: true,
			},
			Alignment: &excelize.Alignment{
				Vertical: "center",
				Indent:   1,
			},
			Border: borderAll,
		})
		centerStyle, _ := file.NewStyle(&excelize.Style{
			CustomNumFmt: &fmtCode,
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
				Indent:     1,
			},
			Border: borderAll,
		})
		borderStyle, _ := file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Vertical: "center",
				Indent:   1,
			},
			Border: borderAll,
		})
		borderHorizontalStyle, _ := file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Vertical: "center",
				Indent:   1,
			},
			Border: borderHorizontal,
		})
		row := 1

		file.SetColWidth(sheet1, "A", "A", 20)
		file.SetColWidth(sheet1, "B", "B", 30)
		file.SetColWidth(sheet1, "C", "D", 20)
		file.SetColWidth(sheet1, "E", "E", 20)
		file.SetColWidth(sheet1, "F", "F", 30)
		file.SetRowHeight(sheet1, row, 30)

		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "A", row), "Buku Besar Piutang")
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "F", row))
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "F", row), titleStyle)

		file.SetRowHeight(sheet1, row, 30)
		row++
		file.SetRowHeight(sheet1, row, 30)
		row++
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "A", row), "Nama ")
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "B", row), data.Contact.Name)
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "C", row), fmt.Sprintf("%s%d", "D", row))
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "E", row), "Kode ")
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "F", row), data.Contact.Code)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "A", row), heroStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "E", row), fmt.Sprintf("%s%d", "E", row), heroStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "B", row), heroStyleNormal)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "F", row), fmt.Sprintf("%s%d", "F", row), heroStyleNormal)
		file.SetRowHeight(sheet1, row, 30)
		row++

		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "A", row), "Alamat ")
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "B", row), data.Contact.Address)
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "C", row), fmt.Sprintf("%s%d", "D", row))
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "E", row), "Telp. ")
		if data.Contact.Phone != nil {
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "F", row), *data.Contact.Phone)
		}
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "A", row), heroStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "E", row), fmt.Sprintf("%s%d", "E", row), heroStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "B", row), heroStyleNormal)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "F", row), fmt.Sprintf("%s%d", "F", row), heroStyleNormal)
		file.SetRowHeight(sheet1, row, 30)
		row++

		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "A", row), "Batas Kredit ")
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "B", row), data.Contact.DebtLimit)
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "C", row), fmt.Sprintf("%s%d", "D", row))
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "E", row), "Sisa Batas Kredit ")
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "F", row), data.Contact.DebtLimitRemain)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "A", row), heroStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "E", row), fmt.Sprintf("%s%d", "E", row), heroStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "B", row), heroStyleWithFormat)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "F", row), fmt.Sprintf("%s%d", "F", row), heroStyleWithFormat)
		file.SetRowHeight(sheet1, row, 30)
		row++
		file.SetRowHeight(sheet1, row, 30)
		row++

		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "A", row), "Tgl")
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "A", row+1))
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "B", row), "Keterangan")
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "B", row+1))
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "C", row), "Ref")
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "C", row), fmt.Sprintf("%s%d", "C", row+1))
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "D", row), "Mutasi")
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "D", row), fmt.Sprintf("%s%d", "E", row))
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "F", row), "Saldo")
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "F", row), fmt.Sprintf("%s%d", "F", row+1))
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "F", row), headerStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "A", row), headerStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "B", row), headerStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "C", row), fmt.Sprintf("%s%d", "C", row), headerStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "D", row), fmt.Sprintf("%s%d", "D", row), headerStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "E", row), fmt.Sprintf("%s%d", "E", row), headerStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "F", row), fmt.Sprintf("%s%d", "F", row), headerStyle)

		file.SetRowHeight(sheet1, row, 30)
		row++
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "D", row), "Debit")
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "E", row), "Kredit")
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "D", row), fmt.Sprintf("%s%d", "E", row), headerStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "C", row), borderHorizontalStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "F", row), fmt.Sprintf("%s%d", "F", row), borderHorizontalStyle)
		file.SetRowHeight(sheet1, row, 30)
		row++

		if data.TotalBalanceBefore > 0 {
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "B", row), "Saldo")
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "B", row), boldLeftStyleFormat)
			file.MergeCell(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "C", row))
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "D", row), data.TotalDebitBefore)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "E", row), data.TotalCreditBefore)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "F", row), data.TotalBalanceBefore)
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "D", row), fmt.Sprintf("%s%d", "F", row), centerStyle)
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "C", row), fmt.Sprintf("%s%d", "C", row), borderStyle)
			file.SetRowHeight(sheet1, row, 30)
			row++
		}

		for _, v := range data.Ledgers {
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "A", row), v.Date.Format("02-01-2006"))
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "B", row), v.Description)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "C", row), v.Ref)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "D", row), v.Debit)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "E", row), v.Credit)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "F", row), v.Balance)
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "C", row), borderStyle)
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "D", row), fmt.Sprintf("%s%d", "F", row), centerStyle)
			file.SetRowHeight(sheet1, row, 30)
			row++

		}

		if data.TotalBalanceAfter > 0 {
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("Saldo %s", input.EndDate.Format("02-01-2006")))
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "B", row), boldLeftStyleFormat)
			file.MergeCell(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "C", row))
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "D", row), data.TotalDebitAfter)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "E", row), data.TotalCreditAfter)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "F", row), data.TotalBalanceAfter)
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "D", row), fmt.Sprintf("%s%d", "F", row), centerStyle)
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "C", row), fmt.Sprintf("%s%d", "C", row), borderStyle)
			file.SetRowHeight(sheet1, row, 30)
			row++
		}

		if data.TotalBalance > 0 {
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "B", row), "Sub Total")
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "B", row), boldLeftStyleFormat)
			file.MergeCell(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "C", row))
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "D", row), data.TotalDebit)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "E", row), data.TotalCredit)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "F", row), data.TotalBalance)
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "D", row), fmt.Sprintf("%s%d", "F", row), centerStyle)
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "C", row), fmt.Sprintf("%s%d", "C", row), borderStyle)
			file.SetRowHeight(sheet1, row, 30)
			row++
		}

		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "B", row), "Total")
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "B", row), boldLeftStyleFormat)
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "C", row))
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "D", row), data.GrandTotalDebit)
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "E", row), data.GrandTotalCredit)
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "F", row), data.GrandTotalBalance)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "D", row), fmt.Sprintf("%s%d", "F", row), centerStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "C", row), fmt.Sprintf("%s%d", "C", row), borderStyle)
		file.SetRowHeight(sheet1, row, 30)
		row++

		fmt.Println(
			headerStyle,
			boldStyle,
			boldStyleFormat,
			centerStyle,
		)

		var buf bytes.Buffer
		if err := file.Write(&buf); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write XLSX file"})
			return
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename=report.xlsx")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account Receivable Ledger report retrieved successfully", "data": data})
}
func (r *ReportHandler) GetAccountPayabledgerHandler(c *gin.Context) {
	var input struct {
		Title      string    `json:"title,omitempty" form:"title"`
		StartDate  time.Time `json:"start_date,omitempty" form:"start_date"`
		EndDate    time.Time `json:"end_date,omitempty" form:"end_date"`
		ContactID  string    `json:"contact_id,omitempty" example:"contact_id"`
		IsDownload bool      `json:"is_download,omitempty" form:"is_download"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	companyID := c.MustGet("companyID").(string)
	// TODO: add date filter
	r.financeSrv.ReportService.SetContactService(r.contactSev)
	data, err := r.financeSrv.ReportService.GetAccountPayableLedger(companyID, input.ContactID, input.StartDate, input.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if input.IsDownload {
		fmtCode := `_(#,##0_);_(\(#,##0\);_("-"??_);_(@_)`
		file := excelize.NewFile()
		sheet1 := file.GetSheetName(0)

		titleStyle, _ := file.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold: true,
				Size: 20,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
				Indent:     10,
			},
			Border: borderAll,
		})
		headerStyle, _ := file.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold: true,
				Size: 12,
			},

			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
				Indent:     1,
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#DCE6F1"}, // Soft blue
				Pattern: 1,
			},
			Border: borderAll,
		})
		heroStyle, _ := file.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold: true,
				Size: 12,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "left",
				Vertical:   "center",
				Indent:     1,
			},
			Border: borderAll,
		})
		heroStyleNormal, _ := file.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Size: 12,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "left",
				Vertical:   "center",
				Indent:     1,
			},
			Border: borderAll,
		})
		heroStyleWithFormat, _ := file.NewStyle(&excelize.Style{
			CustomNumFmt: &fmtCode,
			Font: &excelize.Font{
				Size: 12,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "left",
				Vertical:   "center",
				Indent:     1,
			},
			Border: borderAll,
		})
		boldStyle, _ := file.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold: true,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
				Indent:     1,
			},
			Border: borderAll,
		})
		boldStyleFormat, _ := file.NewStyle(&excelize.Style{
			CustomNumFmt: &fmtCode,
			Font: &excelize.Font{
				Bold: true,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
				Indent:     1,
			},
			Border: borderAll,
		})
		boldLeftStyleFormat, _ := file.NewStyle(&excelize.Style{
			CustomNumFmt: &fmtCode,
			Font: &excelize.Font{
				Bold: true,
			},
			Alignment: &excelize.Alignment{
				Vertical: "center",
				Indent:   1,
			},
			Border: borderAll,
		})
		centerStyle, _ := file.NewStyle(&excelize.Style{
			CustomNumFmt: &fmtCode,
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
				Indent:     1,
			},
			Border: borderAll,
		})
		borderStyle, _ := file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Vertical: "center",
				Indent:   1,
			},
			Border: borderAll,
		})
		borderHorizontalStyle, _ := file.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Vertical: "center",
				Indent:   1,
			},
			Border: borderHorizontal,
		})
		row := 1

		file.SetColWidth(sheet1, "A", "A", 20)
		file.SetColWidth(sheet1, "B", "B", 30)
		file.SetColWidth(sheet1, "C", "D", 20)
		file.SetColWidth(sheet1, "E", "E", 20)
		file.SetColWidth(sheet1, "F", "F", 30)
		file.SetRowHeight(sheet1, row, 30)

		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "A", row), "Buku Besar Hutang")
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "F", row))
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "F", row), titleStyle)

		file.SetRowHeight(sheet1, row, 30)
		row++
		file.SetRowHeight(sheet1, row, 30)
		row++
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "A", row), "Nama ")
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "B", row), data.Contact.Name)
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "C", row), fmt.Sprintf("%s%d", "D", row))
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "E", row), "Kode ")
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "F", row), data.Contact.Code)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "A", row), heroStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "E", row), fmt.Sprintf("%s%d", "E", row), heroStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "B", row), heroStyleNormal)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "F", row), fmt.Sprintf("%s%d", "F", row), heroStyleNormal)
		file.SetRowHeight(sheet1, row, 30)
		row++

		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "A", row), "Alamat ")
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "B", row), data.Contact.Address)
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "C", row), fmt.Sprintf("%s%d", "D", row))
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "E", row), "Telp. ")
		if data.Contact.Phone != nil {
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "F", row), *data.Contact.Phone)
		}
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "A", row), heroStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "E", row), fmt.Sprintf("%s%d", "E", row), heroStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "B", row), heroStyleNormal)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "F", row), fmt.Sprintf("%s%d", "F", row), heroStyleNormal)
		file.SetRowHeight(sheet1, row, 30)
		row++

		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "A", row), "Batas Kredit ")
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "B", row), data.Contact.ReceivablesLimit)
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "C", row), fmt.Sprintf("%s%d", "D", row))
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "E", row), "Sisa Batas Kredit ")
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "F", row), data.Contact.ReceivablesLimitRemain)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "A", row), heroStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "E", row), fmt.Sprintf("%s%d", "E", row), heroStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "B", row), heroStyleWithFormat)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "F", row), fmt.Sprintf("%s%d", "F", row), heroStyleWithFormat)
		file.SetRowHeight(sheet1, row, 30)
		row++
		file.SetRowHeight(sheet1, row, 30)
		row++

		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "A", row), "Tgl")
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "A", row+1))
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "B", row), "Keterangan")
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "B", row+1))
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "C", row), "Ref")
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "C", row), fmt.Sprintf("%s%d", "C", row+1))
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "D", row), "Mutasi")
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "D", row), fmt.Sprintf("%s%d", "E", row))
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "F", row), "Saldo")
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "F", row), fmt.Sprintf("%s%d", "F", row+1))
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "F", row), headerStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "A", row), headerStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "B", row), headerStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "C", row), fmt.Sprintf("%s%d", "C", row), headerStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "D", row), fmt.Sprintf("%s%d", "D", row), headerStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "E", row), fmt.Sprintf("%s%d", "E", row), headerStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "F", row), fmt.Sprintf("%s%d", "F", row), headerStyle)

		file.SetRowHeight(sheet1, row, 30)
		row++
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "D", row), "Debit")
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "E", row), "Kredit")
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "D", row), fmt.Sprintf("%s%d", "E", row), headerStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "C", row), borderHorizontalStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "F", row), fmt.Sprintf("%s%d", "F", row), borderHorizontalStyle)
		file.SetRowHeight(sheet1, row, 30)
		row++

		if data.TotalBalanceBefore > 0 {
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "B", row), "Saldo")
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "B", row), boldLeftStyleFormat)
			file.MergeCell(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "C", row))
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "D", row), data.TotalDebitBefore)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "E", row), data.TotalCreditBefore)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "F", row), data.TotalBalanceBefore)
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "D", row), fmt.Sprintf("%s%d", "F", row), centerStyle)
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "C", row), fmt.Sprintf("%s%d", "C", row), borderStyle)
			file.SetRowHeight(sheet1, row, 30)
			row++
		}

		for _, v := range data.Ledgers {
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "A", row), v.Date.Format("02-01-2006"))
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "B", row), v.Description)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "C", row), v.Ref)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "D", row), v.Debit)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "E", row), v.Credit)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "F", row), v.Balance)
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "A", row), fmt.Sprintf("%s%d", "C", row), borderStyle)
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "D", row), fmt.Sprintf("%s%d", "F", row), centerStyle)
			file.SetRowHeight(sheet1, row, 30)
			row++

		}

		if data.TotalBalanceAfter > 0 {
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("Saldo %s", input.EndDate.Format("02-01-2006")))
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "B", row), boldLeftStyleFormat)
			file.MergeCell(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "C", row))
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "D", row), data.TotalDebitAfter)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "E", row), data.TotalCreditAfter)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "F", row), data.TotalBalanceAfter)
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "D", row), fmt.Sprintf("%s%d", "F", row), centerStyle)
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "C", row), fmt.Sprintf("%s%d", "C", row), borderStyle)
			file.SetRowHeight(sheet1, row, 30)
			row++
		}

		if data.TotalBalance > 0 {
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "B", row), "Sub Total")
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "B", row), boldLeftStyleFormat)
			file.MergeCell(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "C", row))
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "D", row), data.TotalDebit)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "E", row), data.TotalCredit)
			file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "F", row), data.TotalBalance)
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "D", row), fmt.Sprintf("%s%d", "F", row), centerStyle)
			file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "C", row), fmt.Sprintf("%s%d", "C", row), borderStyle)
			file.SetRowHeight(sheet1, row, 30)
			row++
		}

		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "B", row), "Total")
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "B", row), boldLeftStyleFormat)
		file.MergeCell(sheet1, fmt.Sprintf("%s%d", "B", row), fmt.Sprintf("%s%d", "C", row))
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "D", row), data.GrandTotalDebit)
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "E", row), data.GrandTotalCredit)
		file.SetCellValue(sheet1, fmt.Sprintf("%s%d", "F", row), data.GrandTotalBalance)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "D", row), fmt.Sprintf("%s%d", "F", row), centerStyle)
		file.SetCellStyle(sheet1, fmt.Sprintf("%s%d", "C", row), fmt.Sprintf("%s%d", "C", row), borderStyle)
		file.SetRowHeight(sheet1, row, 30)
		row++

		fmt.Println(
			headerStyle,
			boldStyle,
			boldStyleFormat,
			centerStyle,
		)

		var buf bytes.Buffer
		if err := file.Write(&buf); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write XLSX file"})
			return
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename=report.xlsx")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account Payable Ledger report retrieved successfully", "data": data})
}
