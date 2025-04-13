package helper

import (
	"time"

	"github.com/AMETORY/ametory-erp-modules/shared/models"
)

func ValidClosingBook(closingBook *models.ClosingBook, date time.Time) bool {
	if closingBook == nil {
		return true
	}
	return !(closingBook.StartDate.Before(date) && closingBook.EndDate.After(date))
}
