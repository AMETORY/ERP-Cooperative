package objects

import (
	"time"

	"github.com/AMETORY/ametory-erp-modules/shared/models"
)

type SalesRequest struct {
	SalesNumber  string           `json:"sales_number" binding:"required"`
	Code         string           `json:"code"`
	Description  string           `json:"description"`
	Notes        string           `json:"notes"`
	Status       string           `json:"status"`
	SalesDate    time.Time        `json:"sales_date" binding:"required"`
	DueDate      time.Time        `json:"due_date"`
	PaymentTerms string           `json:"payment_terms"`
	ContactID    *string          `json:"contact_id" binding:"required"`
	Type         models.SalesType `json:"type"`
}
