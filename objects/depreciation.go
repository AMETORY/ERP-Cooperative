package objects

import (
	"time"

	"github.com/AMETORY/ametory-erp-modules/shared/models"
)

type DepreciationRequest struct {
	Date                             time.Time                      `json:"date" binding:"required"`
	AccountCurrentAssetID            string                         `json:"account_current_asset_id" binding:"required"`
	AccountFixedAssetID              string                         `json:"account_fixed_asset_id" binding:"required"`
	AccountDepreciationID            string                         `json:"account_depreciation_id" binding:"required"`
	AccountAccumulatedDepreciationID string                         `json:"account_accumulated_depreciation_id" binding:"required"`
	DepreciationCosts                []models.DepreciationCostModel `json:"depreciation_costs" binding:"required"`
	DepreciationMethod               string                         `json:"depreciation_method" binding:"required"`
	IsMonthly                        bool                           `json:"is_monthly"`
}
