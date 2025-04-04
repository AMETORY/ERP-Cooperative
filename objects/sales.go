package objects

import (
	"time"

	"github.com/AMETORY/ametory-erp-modules/shared/models"
)

type SalesRequest struct {
	SalesNumber      string                  `json:"sales_number" binding:"required"`
	Code             string                  `json:"code"`
	Description      string                  `json:"description"`
	Notes            string                  `json:"notes"`
	Status           string                  `json:"status"`
	SalesDate        time.Time               `json:"sales_date" binding:"required"`
	DueDate          *time.Time              `json:"due_date"`
	PaymentTerms     string                  `json:"payment_terms"`
	ContactID        *string                 `json:"contact_id" binding:"required"`
	Type             models.SalesType        `json:"type"`
	DocumentType     models.SalesDocType     `json:"document_type"`
	Items            []models.SalesItemModel `json:"items"`
	RefID            *string                 `json:"ref_id,omitempty"`
	RefType          *string                 `json:"ref_type,omitempty"`
	SecondaryRefID   *string                 `json:"secondary_ref_id,omitempty"`
	SecondaryRefType *string                 `json:"secondary_ref_type,omitempty"`
	PaymentTermsCode string                  `json:"payment_terms_code"`
	TermCondition    string                  `json:"term_condition"`
	DeliveryID       *string                 `json:"delivery_id"`
	DeliveryData     string                  `gorm:"type:json" json:"delivery_data"`
}
