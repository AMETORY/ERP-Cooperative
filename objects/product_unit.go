package objects

type ProductUnitRequest struct {
	UnitID    string  `json:"unit_id" binding:"required"`
	Value     float64 `json:"value" binding:"required"`
	IsDefault bool    `json:"is_default"`
}
