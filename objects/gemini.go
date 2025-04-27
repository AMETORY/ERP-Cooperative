package objects

type GeminiResponse struct {
	Response  string `json:"response"`
	Type      string `json:"type"`
	Command   string `json:"command"`
	Params    any    `json:"params"`
	CompanyID string `json:"company_id"`
	UserID    string `json:"user_id"`
}
