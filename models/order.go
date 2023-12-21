package models

type Order struct {
	OrderNum  string  `json:"number"`
	Status    string  `json:"status"`
	Accrual   float64 `json:"accrual,omitempty"`
	Sum       float64 `json:"sum,omitempty"`
	UUID      string  `json:"uuid,omitempty"`
	CreatedAt string  `json:"uploaded_at"`
}
