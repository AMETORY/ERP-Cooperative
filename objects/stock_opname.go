package objects

import "time"

type StockOpnameRequest struct {
	Date        time.Time `json:"date" binding:"required"`
	InventoryID *string   `json:"inventory_id"`
	ExpenseID   *string   `json:"expense_id"`
	RevenueID   *string   `json:"revenue_id"`
	Notes       string    `json:"notes"`
}
