package pos

import (
	"time"

	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"github.com/gin-gonic/gin"
)

func (s *PosHandler) GetDashboardSummaryHandler(c *gin.Context) {
	merchantID := c.MustGet("merchantID").(string)
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

	var countSales int64
	var totalSalesOrder float64
	var totalSales float64

	err := s.ctx.DB.Model(&models.MerchantOrder{}).
		Where("cashier_id = ?", userID).
		Where("created_at  between ? AND ?", input.StartDate, input.EndDate).
		Where("order_status = ?", "PAID").
		Where("merchant_id = ?", merchantID).
		Count(&countSales).Error
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	err = s.ctx.DB.Model(&models.MerchantOrder{}).
		Where("cashier_id = ?", userID).
		Where("created_at  between ? AND ?", input.StartDate, input.EndDate).
		Where("order_status = ?", "PAID").
		Select("COALESCE(sum(total), 0) as total").
		Where("merchant_id = ?", merchantID).
		Scan(&totalSalesOrder).Error
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	err = s.ctx.DB.Model(&models.MerchantOrder{}).
		Where("cashier_id = ?", userID).
		Where("created_at  between ? AND ?", input.StartDate, input.EndDate).
		Where("order_status = ?", "PAID").
		Select("COALESCE(sum(total), 0) as total").
		Where("merchant_id = ?", merchantID).
		Scan(&totalSales).Error
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	type summary struct {
		Day       string  `json:"day"`
		Count     int64   `json:"count"`
		Total     float64 `json:"total"`
		DayOfWeek int     `json:"day_of_week"`
	}

	var summaryPerDay []summary

	start := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, -int(time.Monday-time.Now().Weekday()))

	for i := range 7 {
		day := start.AddDate(0, 0, i)
		var count int64
		var total float64

		err = s.ctx.DB.Model(&models.MerchantOrder{}).
			Where("cashier_id = ?", userID).
			Where("created_at  between ? AND ?", day, day.AddDate(0, 0, 1).Add(-time.Nanosecond)).
			Where("order_status = ?", "PAID").
			Where("merchant_id = ?", merchantID).
			Count(&count).Error
		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}
		err = s.ctx.DB.Model(&models.MerchantOrder{}).
			Where("cashier_id = ?", userID).
			Where("created_at  between ? AND ?", day, day.AddDate(0, 0, 1).Add(-time.Nanosecond)).
			Where("order_status = ?", "PAID").
			Select("COALESCE(sum(total), 0) as total").
			Where("merchant_id = ?", merchantID).
			Scan(&total).Error
		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}

		summaryPerDay = append(summaryPerDay, summary{
			DayOfWeek: int(day.Weekday()),
			Day:       day.Format("Mon"),
			Count:     count,
			Total:     total,
		})
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data": gin.H{
			"count_sales":       countSales,
			"total_sales_order": totalSalesOrder,
			"total_sales":       totalSales,
			"summary_per_day":   summaryPerDay,
		},
	})

}
