package models

type Balance struct {
	Current  float32 `json:"current"`
	Withdraw float32 `json:"withdrawn"`
}

type Drawall struct {
	Order       string  `json:"order"`
	Sum         float32 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
