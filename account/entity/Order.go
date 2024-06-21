package entity

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	OID       uuid.UUID
	TID       uuid.UUID
	TPrice    float64
	Status    string
	CreatedAt time.Time
	LastEdit  time.Time
}

const (
	OrderStatusNew        = "New"
	OrderStatusPaid       = "Paid"
	OrderStatusProcessing = "Processing"
	OrderStatusDone       = "Done"
)
