package handlers

import (
	"ametory-cooperative/app_models"
	"fmt"

	"github.com/AMETORY/ametory-erp-modules/shared"
	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"gorm.io/gorm"
)

func GenerateAutoNumber(db *gorm.DB, companyID string, docType string) (string, error) {
	var setting app_models.CustomSettingModel
	db.Where("id = ?", companyID).First(&setting)
	if setting.AutoNumericLength == nil || setting.RandomCharacterLength == nil || setting.RandomNumericLength == nil {
		return "", nil
	}

	data := shared.InvoiceBillSettingModel{
		AutoNumericLength:     *setting.AutoNumericLength,
		RandomNumericLength:   *setting.RandomNumericLength,
		RandomCharacterLength: *setting.RandomCharacterLength,
	}
	nextNumber := ""
	switch docType {
	case string(models.INVOICE):
		if setting.SalesStaticCharacter == nil || setting.SalesFormat == nil {
			return "", nil
		}
		data.StaticCharacter = *setting.SalesStaticCharacter
		data.NumberFormat = *setting.SalesFormat
	case string(models.SALES_ORDER):
		if setting.SalesOrderStaticCharacter == nil || setting.SalesOrderFormat == nil {
			return "", nil
		}
		data.StaticCharacter = *setting.SalesOrderStaticCharacter
		data.NumberFormat = *setting.SalesOrderFormat
	case string(models.SALES_QUOTE):
		if setting.SalesQuoteStaticCharacter == nil || setting.SalesQuoteFormat == nil {
			return "", nil
		}
		data.StaticCharacter = *setting.SalesQuoteStaticCharacter
		data.NumberFormat = *setting.SalesQuoteFormat
	case string(models.DELIVERY):
		if setting.DeliveryStaticCharacter == nil || setting.DeliveryFormat == nil {
			return "", nil
		}
		data.StaticCharacter = *setting.DeliveryStaticCharacter
		data.NumberFormat = *setting.DeliveryFormat
	case string(models.BILL):
		if setting.PurchaseStaticCharacter == nil || setting.PurchaseFormat == nil {
			return "", nil
		}
		data.StaticCharacter = *setting.PurchaseStaticCharacter
		data.NumberFormat = *setting.PurchaseFormat
	case string(models.PURCHASE_ORDER):
		if setting.PurchaseOrderStaticCharacter == nil || setting.PurchaseOrderFormat == nil {
			return "", nil
		}
		data.StaticCharacter = *setting.PurchaseOrderStaticCharacter
		data.NumberFormat = *setting.PurchaseOrderFormat

	}

	switch docType {
	case string(models.INVOICE), string(models.SALES_ORDER), string(models.SALES_QUOTE), string(models.DELIVERY):
		lastDoc := models.SalesModel{}
		if err := db.Where("company_id = ? AND document_type = ?", companyID, docType).Limit(1).Order("updated_at desc").First(&lastDoc).Error; err != nil {
			nextNumber = shared.GenerateInvoiceBillNumber(data, "00")
			fmt.Println("KESINI #1", nextNumber)
		} else {

			nextNumber = shared.ExtractNumber(data, lastDoc.SalesNumber)
			fmt.Println("KESINI #2", nextNumber)
		}
	case string(models.BILL), string(models.PURCHASE_ORDER):
		lastDoc := models.PurchaseOrderModel{}
		if err := db.Where("company_id = ?  AND document_type = ?", companyID, docType).Limit(1).Order("updated_at desc").First(&lastDoc).Error; err != nil {
			nextNumber = shared.GenerateInvoiceBillNumber(data, "00")
			fmt.Println("KESINI #3", nextNumber)
		} else {
			nextNumber = shared.ExtractNumber(data, lastDoc.PurchaseNumber)
			fmt.Println("KESINI #4", nextNumber)
		}
	}

	return nextNumber, nil
}
