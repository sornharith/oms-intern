package model

import "time"

type Order struct {
	OID       int       `db:"o_id" json:"o_id"`
	TID       int       `db:"t_id" json:"t_id"`
	TPrice    float64   `db:"t_price" json:"t_price"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"create_at" json:"create_at"`
	LastEdit  time.Time `db:"last_edit" json:"last_edit"`
}
