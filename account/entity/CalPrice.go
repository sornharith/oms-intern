package entity

import "github.com/google/uuid"

type CalPrice struct {
	TID        uuid.UUID `db:"t_id"`
	TPrice     float64   `db:"t_price"`
	UserSelect string    `db:"user_select"`
	Address    string    `db:"address"`
}

type UserSelectItem struct {
	ProductID int `json:"product_id"`
	Amount    int `json:"amount"`
}
