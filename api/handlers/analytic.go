package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/finance"
	"github.com/AMETORY/ametory-erp-modules/inventory"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type AnalyticHandler struct {
	ctx          *context.ERPContext
	inventorySrv *inventory.InventoryService
	financeSrv   *finance.FinanceService
}

func NewAnalyticHandler(ctx *context.ERPContext) *AnalyticHandler {
	inventorySrv, ok := ctx.InventoryService.(*inventory.InventoryService)
	if !ok {
		panic("inventory service is not found")
	}
	financeSrv, ok := ctx.FinanceService.(*finance.FinanceService)
	if !ok {
		panic("finance service is not found")
	}
	return &AnalyticHandler{
		ctx:          ctx,
		inventorySrv: inventorySrv,
		financeSrv:   financeSrv,
	}
}

func (a *AnalyticHandler) PopularProductHandler(c *gin.Context) {
	data, err := a.inventorySrv.ProductService.GetBestSellingProduct(c.Request, 5, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Popular product retrieved successfully", "data": data})
}

func (a *AnalyticHandler) GetMonthlySalesReportHandler(c *gin.Context) {
	companyID := c.MustGet("companyID").(string)
	year := c.Query("year")

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		yearInt = time.Now().Year()
	}
	data, err := a.financeSrv.ReportService.GetMonthlySalesReport(companyID, yearInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Monthly Sales Report retrieved successfully", "data": data})
}
func (a *AnalyticHandler) GetMonthlySalesPurchaseReportHandler(c *gin.Context) {
	companyID := c.MustGet("companyID").(string)
	year := c.Query("year")

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		yearInt = time.Now().Year()
	}
	sales, err := a.financeSrv.ReportService.GetMonthlySalesReport(companyID, yearInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	purchase, err := a.financeSrv.ReportService.GetMonthlyPurchaseReport(companyID, yearInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	salesPurchase := []struct {
		Month    string  `json:"month"`
		Sales    float64 `json:"sales"`
		Purchase float64 `json:"purchase"`
	}{}

	for i, v := range sales {
		salesPurchase = append(salesPurchase, struct {
			Month    string  `json:"month"`
			Sales    float64 `json:"sales"`
			Purchase float64 `json:"purchase"`
		}{
			Month:    v.MonthName,
			Sales:    v.Total,
			Purchase: purchase[i].Total,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Monthly Sales Report retrieved successfully", "data": gin.H{"data": salesPurchase, "title": "Monthly Sales & Purchase"}})
}
func (a *AnalyticHandler) GetMonthlyPurchaseReportHandler(c *gin.Context) {
	companyID := c.MustGet("companyID").(string)
	year := c.Query("year")

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		yearInt = time.Now().Year()
	}
	data, err := a.financeSrv.ReportService.GetMonthlyPurchaseReport(companyID, yearInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Monthly Purchase Report retrieved successfully", "data": data})
}

func (a *AnalyticHandler) GetSalesTimeRangeHandler(c *gin.Context) {
	companyID := c.MustGet("companyID").(string)

	timeRange := c.DefaultQuery("time_range", "THIS_WEEK")
	docType := c.DefaultQuery("doc_type", "INVOICE")

	data, err := a.financeSrv.ReportService.CalculateSalesByTimeRange(companyID, docType, timeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sales report for time range retrieved successfully", "data": data})
}

func (a *AnalyticHandler) GetPurchaseTimeRangeHandler(c *gin.Context) {
	companyID := c.MustGet("companyID").(string)

	timeRange := c.DefaultQuery("time_range", "THIS_WEEK")
	docType := c.DefaultQuery("doc_type", "BILL")

	data, err := a.financeSrv.ReportService.CalculatePurchaseByTimeRange(companyID, docType, timeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase report for time range retrieved successfully", "data": data})
}

func (a *AnalyticHandler) GetWeeklyPurchaseReportHandler(c *gin.Context) {
	companyID := c.MustGet("companyID").(string)

	year := c.Query("year")
	month := c.Query("month")

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		yearInt = time.Now().Year()
	}
	monthInt, err := strconv.Atoi(month)
	if err != nil {
		monthInt = int(time.Now().Month())
	}

	data, err := a.financeSrv.ReportService.GetWeeklyPurchaseReport(companyID, yearInt, monthInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Weekly Purchase report retrieved successfully", "data": data})
}

func (a *AnalyticHandler) GetWeeklySalesReportHandler(c *gin.Context) {
	companyID := c.MustGet("companyID").(string)

	year := c.Query("year")
	month := c.Query("month")

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		yearInt = time.Now().Year()
	}
	monthInt, err := strconv.Atoi(month)
	if err != nil {
		monthInt = int(time.Now().Month())
	}

	data, err := a.financeSrv.ReportService.GetWeeklySalesReport(companyID, yearInt, monthInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Weekly Sales report retrieved successfully", "data": data})
}

func (a *AnalyticHandler) GetWeeklySalesPurchaseReportHandler(c *gin.Context) {
	companyID := c.MustGet("companyID").(string)

	year := c.Query("year")
	month := c.Query("month")

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		yearInt = time.Now().Year()
	}
	monthInt, err := strconv.Atoi(month)
	if err != nil {
		monthInt = int(time.Now().Month())
	}

	sales, err := a.financeSrv.ReportService.GetWeeklySalesReport(companyID, yearInt, monthInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	purchase, err := a.financeSrv.ReportService.GetWeeklyPurchaseReport(companyID, yearInt, monthInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	salesPurchase := []struct {
		Week     string  `json:"week"`
		Sales    float64 `json:"sales"`
		Purchase float64 `json:"purchase"`
	}{}

	for i, v := range sales {
		salesPurchase = append(salesPurchase, struct {
			Week     string  `json:"week"`
			Sales    float64 `json:"sales"`
			Purchase float64 `json:"purchase"`
		}{
			Week:     v.WeekName,
			Sales:    v.Total,
			Purchase: purchase[i].Total,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Weekly Sales & Purchase Report retrieved successfully", "data": gin.H{"data": salesPurchase, "title": "Weekly Sales & Purchase"}})
}

func (a *AnalyticHandler) GetNetWorthHandler(c *gin.Context) {
	companyID := c.MustGet("companyID").(string)
	var start, end *time.Time
	if c.Request.Header.Get("start-date") != "" {
		startDate, err := time.Parse(time.RFC3339, c.Request.Header.Get("start-date"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			c.Abort()
			return
		}
		start = &startDate
	}
	if c.Request.Header.Get("end-date") != "" {
		endDate, err := time.Parse(time.RFC3339, c.Request.Header.Get("end-date"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			c.Abort()
			return
		}
		end = &endDate
	}
	report, err := a.financeSrv.ReportService.GenerateProfitLossReport(models.GeneralReport{
		StartDate: *start,
		EndDate:   *end,
		CompanyID: companyID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Net Worth Report retrieved successfully", "data": report})
}

func (a *AnalyticHandler) GetSumCashBankHandler(c *gin.Context) {
	companyID := c.MustGet("companyID").(string)
	data, err := a.financeSrv.ReportService.GetSumCashBank(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sum of cash and bank retrieved successfully", "data": data})
}
