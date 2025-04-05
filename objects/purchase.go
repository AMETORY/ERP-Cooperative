package objects

import (
	"time"

	"github.com/AMETORY/ametory-erp-modules/shared/models"
)

type PurchaseRequest struct {
	PurchaseNumber   string                          `json:"purchase_number" binding:"required"`
	Code             string                          `json:"code"`
	Description      string                          `json:"description"`
	Notes            string                          `json:"notes"`
	Status           string                          `json:"status"`
	PurchaseDate     time.Time                       `json:"purchase_date" binding:"required"`
	DueDate          *time.Time                      `json:"due_date"`
	PaymentTerms     string                          `json:"payment_terms"`
	ContactID        *string                         `json:"contact_id" binding:"required"`
	Type             models.PurchaseType             `json:"type"`
	DocumentType     models.PurchaseDocType          `json:"document_type"`
	Items            []models.PurchaseOrderItemModel `json:"items"`
	RefID            *string                         `json:"ref_id,omitempty"`
	RefType          *string                         `json:"ref_type,omitempty"`
	SecondaryRefID   *string                         `json:"secondary_ref_id,omitempty"`
	SecondaryRefType *string                         `json:"secondary_ref_type,omitempty"`
	PaymentTermsCode string                          `json:"payment_terms_code"`
	TermCondition    string                          `json:"term_condition"`
	DeliveryID       *string                         `json:"delivery_id"`
	DeliveryData     string                          `gorm:"type:json" json:"delivery_data"`
}
