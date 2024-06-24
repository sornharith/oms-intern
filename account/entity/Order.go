package entity

import (
	"github.com/google/uuid"
	service "memrizr/account/service/model"
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

func (o *Order) IsValidStatus(status service.OrderStatus) bool {
	if o.Status == OrderStatusNew && status != OrderStatusPaid {
		return false
	}
	return status == OrderStatusPaid ||
		status == OrderStatusProcessing ||
		status == OrderStatusDone
}
