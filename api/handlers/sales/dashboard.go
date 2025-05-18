package sales

import (
	"fmt"
	"time"

	"github.com/AMETORY/ametory-erp-modules/context"
	"github.com/AMETORY/ametory-erp-modules/finance"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

type SalesDashboardHandler struct {
	ctx        *context.ERPContext
	financeSrv *finance.FinanceService
}

func NewSalesDashboardHandler(ctx *context.ERPContext) *SalesDashboardHandler {
	financeSrv, ok := ctx.FinanceService.(*finance.FinanceService)
	if !ok {
		panic("finance service is not found")
	}
	return &SalesDashboardHandler{
		ctx:        ctx,
		financeSrv: financeSrv,
	}
}
func (s *SalesDashboardHandler) GetDashboardSummaryHandler(c *gin.Context) {
	input := struct {
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
		Year      int       `json:"year"`
	}{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	userID := c.MustGet("userID").(string)
	companyID := c.MustGet("companyID").(string)

	var countSales int64
	var totalSalesOrder float64
	var totalSales float64

	err := s.ctx.DB.Model(&models.SalesModel{}).
		Where("document_type = ? AND company_id = ? AND sales_user_id = ?", models.SALES_ORDER, companyID, userID).
		Where("sales_date  between ? AND ?", input.StartDate, input.EndDate).
		Where("status IN (?)", []string{"RELEASED"}).
		Count(&countSales).Error
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	err = s.ctx.DB.Model(&models.SalesModel{}).
		Where("document_type = ? AND company_id = ? AND sales_user_id = ?", models.SALES_ORDER, companyID, userID).
		Where("sales_date  between ? AND ?", input.StartDate, input.EndDate).
		Select("COALESCE(sum(total), 0) as total").
		Where("status IN (?)", []string{"RELEASED"}).
		Scan(&totalSalesOrder).Error
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	err = s.ctx.DB.Model(&models.SalesModel{}).
		Where("document_type = ? AND company_id = ? AND sales_user_id = ?", models.INVOICE, companyID, userID).
		Where("sales_date  between ? AND ?", input.StartDate, input.EndDate).
		Where("status IN (?)", []string{"POSTED"}).
		Select("COALESCE(sum(total), 0) as total").
		Scan(&totalSales).Error
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	type SalesSummaryPerMonth struct {
		Month      string  `json:"month"`
		TotalCount int64   `json:"total_count"`
		Total      float64 `json:"total"`
	}

	fmt.Println("YEAR", input.Year)

	summaryPerMonth := make([]SalesSummaryPerMonth, 12)
	for i := 0; i < 12; i++ {
		month := time.Date(input.Year, time.Month(i+1), 1, 0, 0, 0, 0, input.StartDate.Location())
		endMonth := month.AddDate(0, 1, -1)

		var total float64
		err = s.ctx.DB.Model(&models.SalesModel{}).
			Where("company_id = ? AND sales_user_id = ?", companyID, userID).
			Where("sales_date  between ? AND ?", month, endMonth).
			Where("document_type = ?", models.INVOICE).
			Where("status IN (?)", []string{"POSTED"}).
			Select("COALESCE(sum(total), 0) as total").
			Scan(&total).Error
		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}

		summaryPerMonth[i] = SalesSummaryPerMonth{
			Month: month.Format("January"),
			Total: total,
		}
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data": gin.H{
			"count_sales":       countSales,
			"total_sales_order": totalSalesOrder,
			"total_sales":       totalSales,
			"summary_per_month": summaryPerMonth,
		},
	})

}
