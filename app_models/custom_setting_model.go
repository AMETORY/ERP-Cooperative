package app_models

import "github.com/AMETORY/ametory-erp-modules/shared/models"

type CustomSettingModel struct {
	models.CompanyModel
	GeminiAPIKey          *string                         `gorm:"type:varchar(255)" json:"gemini_api_key,omitempty"`
	WhatsappWebHost       *string                         `gorm:"type:varchar(255)" json:"whatsapp_web_host,omitempty"`
	WhatsappWebMockNumber *string                         `gorm:"type:varchar(255)" json:"whatsapp_web_mock_number,omitempty"`
	WhatsappWebIsMocked   *string                         `gorm:"type:varchar(255)" json:"whatsapp_web_is_mocked,omitempty"`
	CooperativeSetting    *models.CooperativeSettingModel `gorm:"-" json:"cooperative_setting,omitempty"`
}

func (CustomSettingModel) TableName() string {
	return "companies"
}
