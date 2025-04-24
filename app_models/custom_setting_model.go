package app_models

import (
	"encoding/json"
	"time"

	"github.com/AMETORY/ametory-erp-modules/shared/models"
	"gorm.io/gorm"
)

type CustomSettingModel struct {
	models.CompanyModel
	GeminiAPIKey                  *string                         `gorm:"type:varchar(255)" json:"gemini_api_key,omitempty"`
	WhatsappWebHost               *string                         `gorm:"type:varchar(255)" json:"whatsapp_web_host,omitempty"`
	WhatsappWebMockNumber         *string                         `gorm:"type:varchar(255)" json:"whatsapp_web_mock_number,omitempty"`
	WhatsappWebIsMocked           *string                         `gorm:"type:varchar(255)" json:"whatsapp_web_is_mocked,omitempty"`
	CooperativeSetting            *models.CooperativeSettingModel `gorm:"-" json:"cooperative_setting,omitempty"`
	IsPremium                     bool
	PremiumExpiredAt              *time.Time
	SalesStaticCharacter          *string `json:"sales_static_character"`
	SalesOrderStaticCharacter     *string `json:"sales_order_static_character"`
	SalesQuoteStaticCharacter     *string `json:"sales_quote_static_character"`
	SalesReturnStaticCharacter    *string `json:"sales_return_static_character"`
	DeliveryStaticCharacter       *string `json:"delivery_static_character"`
	PurchaseStaticCharacter       *string `json:"purchase_static_character"`
	PurchaseOrderStaticCharacter  *string `json:"purchase_order_static_character"`
	PurchaseReturnStaticCharacter *string `json:"purchase_return_static_character"`
	SalesFormat                   *string `json:"sales_format"`
	SalesOrderFormat              *string `json:"sales_order_format"`
	SalesQuoteFormat              *string `json:"sales_quote_format"`
	SalesReturnFormat             *string `json:"sales_return_format"`
	DeliveryFormat                *string `json:"delivery_format"`
	PurchaseFormat                *string `json:"purchase_format"`
	PurchaseOrderFormat           *string `json:"purchase_order_format"`
	PurchaseReturnFormat          *string `json:"purchase_return_format"`
	AutoNumericLength             *int    `json:"auto_numeric_length"`
	RandomNumericLength           *int    `json:"random_numeric_length"`
	RandomCharacterLength         *int    `json:"random_character_length"`
}

func (CustomSettingModel) TableName() string {
	return "companies"
}

func (u *CustomSettingModel) AfterFind(tx *gorm.DB) error {
	if u.CashflowGroupSettingData != nil {
		if err := json.Unmarshal([]byte(*u.CashflowGroupSettingData), &u.CashflowGroupSetting); err != nil {
			return err
		}
	}
	if u.CashflowGroupSettingData == nil {
		defaultSetting := models.DefaultCasflowGroupSetting()
		if u.IsCooperation {
			defaultSetting = models.CooperationCasflowGroupSetting()
		}
		u.CashflowGroupSetting = &defaultSetting
		b, _ := json.Marshal(defaultSetting)
		settingStr := string(b)
		u.CashflowGroupSettingData = &settingStr
	}
	if u.RandomNumericLength == nil {
		u.RandomNumericLength = new(int)
		*u.RandomNumericLength = 4
	}
	if u.RandomCharacterLength == nil {
		u.RandomCharacterLength = new(int)
		*u.RandomCharacterLength = 2
	}

	if u.AutoNumericLength == nil {
		u.AutoNumericLength = new(int)
		*u.AutoNumericLength = 4
	}

	if u.SalesStaticCharacter == nil {
		u.SalesStaticCharacter = new(string)
		*u.SalesStaticCharacter = "SALES"
	}
	if u.SalesOrderStaticCharacter == nil {
		u.SalesOrderStaticCharacter = new(string)
		*u.SalesOrderStaticCharacter = "SO"
	}
	if u.SalesQuoteStaticCharacter == nil {
		u.SalesQuoteStaticCharacter = new(string)
		*u.SalesQuoteStaticCharacter = "QUOTE"
	}
	if u.SalesReturnStaticCharacter == nil {
		u.SalesReturnStaticCharacter = new(string)
		*u.SalesReturnStaticCharacter = "SLS-RTN"
	}
	if u.DeliveryStaticCharacter == nil {
		u.DeliveryStaticCharacter = new(string)
		*u.DeliveryStaticCharacter = "DO"
	}
	if u.PurchaseStaticCharacter == nil {
		u.PurchaseStaticCharacter = new(string)
		*u.PurchaseStaticCharacter = "PUR"
	}
	if u.PurchaseOrderStaticCharacter == nil {
		u.PurchaseOrderStaticCharacter = new(string)
		*u.PurchaseOrderStaticCharacter = "PO"
	}
	if u.PurchaseReturnStaticCharacter == nil {
		u.PurchaseReturnStaticCharacter = new(string)
		*u.PurchaseReturnStaticCharacter = "PUR-RTN"
	}

	if u.SalesFormat == nil {
		u.SalesFormat = new(string)
		*u.SalesFormat = "{static-character}-{auto-numeric}/{month-roman}/{year-yyyy}"
	}
	if u.SalesOrderFormat == nil {
		u.SalesOrderFormat = new(string)
		*u.SalesOrderFormat = "{static-character}-{auto-numeric}/{month-roman}/{year-yyyy}"
	}
	if u.SalesQuoteFormat == nil {
		u.SalesQuoteFormat = new(string)
		*u.SalesQuoteFormat = "{static-character}-{auto-numeric}/{month-roman}/{year-yyyy}"
	}
	if u.SalesReturnFormat == nil {
		u.SalesReturnFormat = new(string)
		*u.SalesReturnFormat = "{static-character}-{auto-numeric}/{month-roman}/{year-yyyy}"
	}
	if u.DeliveryFormat == nil {
		u.DeliveryFormat = new(string)
		*u.DeliveryFormat = "{static-character}-{auto-numeric}/{month-roman}/{year-yyyy}"
	}
	if u.PurchaseFormat == nil {
		u.PurchaseFormat = new(string)
		*u.PurchaseFormat = "{static-character}-{auto-numeric}/{month-roman}/{year-yyyy}"
	}
	if u.PurchaseOrderFormat == nil {
		u.PurchaseOrderFormat = new(string)
		*u.PurchaseOrderFormat = "{static-character}-{auto-numeric}/{month-roman}/{year-yyyy}"
	}
	if u.PurchaseReturnFormat == nil {
		u.PurchaseReturnFormat = new(string)
		*u.PurchaseReturnFormat = "{static-character}-{auto-numeric}/{month-roman}/{year-yyyy}"
	}

	return tx.Save(u).Error
}
