package objects

import "time"

type TransactionRequest struct {
	SourceID         string    `json:"source_id"`
	AccountID        *string   `json:"account_id"`
	DestinationID    string    `json:"destination_id"`
	Amount           float64   `json:"amount" binding:"required"`
	Credit           float64   `json:"credit"`
	Debit            float64   `json:"debit"`
	Description      string    `json:"description" binding:"required"`
	Note             string    `json:"note"`
	Date             time.Time `json:"date" binding:"required"`
	IsIncome         bool      `json:"is_income"`
	IsExpense        bool      `json:"is_expense"`
	IsEquity         bool      `json:"is_equity"`
	IsTransfer       bool      `json:"is_transfer"`
	IsOpeningBalance bool      `json:"is_opening_balance"`
}
type TransactionUpdateRequest struct {
	Amount      float64   `json:"amount" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
}
