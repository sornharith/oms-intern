package entity

import "github.com/google/uuid"

type CalPrice struct {
	TID        uuid.UUID
	TPrice     float64
	UserSelect string
	Address    string
}

type UserSelectItem struct {
	ProductID int `json:"product_id"`
	Amount    int `json:"amount"`
}
