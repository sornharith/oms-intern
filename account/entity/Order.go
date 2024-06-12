package entity

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	OID       uuid.UUID `db:"o_id" json:"o_id"`
	TID       uuid.UUID `db:"t_id" json:"t_id"`
	TPrice    float64   `db:"t_price" json:"t_price"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"create_at" json:"create_at"`
	LastEdit  time.Time `db:"last_edit" json:"last_edit"`
}

const (
	OrderStatusNew        = "New"
	OrderStatusPaid       = "Paid"
	OrderStatusProcessing = "Processing"
	OrderStatusDone       = "Done"
)
