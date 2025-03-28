package objects

import "github.com/AMETORY/ametory-erp-modules/shared/models"

type CompanyRequest struct {
	Name              string                `json:"name" binding:"required"`
	Address           string                `json:"address" binding:"required"`
	Email             string                `json:"email" binding:"required"`
	Phone             string                `json:"phone" binding:"required"`
	SectorID          string                `json:"sector_id" binding:"required"`
	CompanyCategoryID string                `json:"company_category_id" binding:"required"`
	IsIslamic         bool                  `json:"is_islamic"`
	ProvinceID        string                `json:"province_id" binding:"required"`
	RegencyID         string                `json:"regency_id" binding:"required"`
	DistrictID        string                `json:"district_id" binding:"required"`
	VillageID         string                `json:"village_id" binding:"required"`
	ZipCode           string                `json:"zip_code" binding:"required"`
	IsCooperation     bool                  `json:"is_cooperation"`
	Accounts          []models.AccountModel `json:"accounts,omitempty"`
}
