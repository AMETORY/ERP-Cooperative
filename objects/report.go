package objects

import "time"

type ReportRequest struct {
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	ProductIDs  []string  `json:"product_ids"`
	CustomerIDs []string  `json:"customer_ids"`
	View        string    `json:"view"`
	IsDownload  bool      `json:"is_download"`
}
