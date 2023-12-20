package models

type Order struct {
	OrderNum  string `json:"number"`
	Status    string `json:"status"`
	Accrual   int    `json:"accrual,omitempty"`
	UUID      string `json:"uuid,omitempty"`
	CreatedAt string `json:"uploaded_at"`
}
