package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	service "memrizr/account/service/model"
)

func TestOrder_IsValidStatus(t *testing.T) {
	tests := []struct {
		initialStatus string
		newStatus     string
		expected      bool
	}{
		{OrderStatusNew, OrderStatusPaid, true},
		{OrderStatusPaid, OrderStatusProcessing, true},
		{OrderStatusProcessing, OrderStatusDone, true},
		{OrderStatusPaid, OrderStatusDone, true},
		{OrderStatusNew, OrderStatusProcessing, false},
		{OrderStatusNew, OrderStatusDone, false},
	}

	for _, test := range tests {
		order := Order{
			OID:       uuid.New(),
			TID:       uuid.New(),
			TPrice:    100.0,
			Status:    test.initialStatus,
			CreatedAt: time.Now(),
			LastEdit:  time.Now(),
		}

		status := service.OrderStatus(test.newStatus)
		result := order.IsValidStatus(status)

		if result != test.expected {
			t.Errorf("IsValidStatus(%v) = %v; want %v", status, result, test.expected)
		}
	}
}
