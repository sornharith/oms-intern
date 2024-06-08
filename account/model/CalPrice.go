package model

type CalPrice struct {
	TID        int     `db:"t_id" json:"t_id"`
	TPrice     float64 `db:"t_price" json:"t_price"`
	UserSelect string  `db:"user_select" json:"user_select"`
	Address    string  `db:"address" json:"address"`
}
