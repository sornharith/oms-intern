package responsemodel

import "github.com/google/uuid"

type CalPrice struct {
	TID uuid.UUID `json:"t_id"`
	TPrice float64 `json:"t_price"`
	UserSelect string `json:"user_select"`
	Address string `json:"address"`
}