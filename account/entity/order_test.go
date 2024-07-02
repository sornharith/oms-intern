package entity

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	service "memrizr/account/service/model"
	"testing"
	"time"
)

const (
	MockOrderStatusPaid       = "Paid"
	MockOrderStatusProcessing = "Processing"
	MockOrderStatusDone       = "Done"
	MockOrderStatusInvalid    = "Invalid"
)

func TestOrder_IsValidStatus(t *testing.T) {
	t.Run("Test Valid Status Transition from New to Paid", func(t *testing.T) {
		order := &Order{
			OID:       uuid.New(),
			TID:       uuid.New(),
			TPrice:    100.0,
			Status:    OrderStatusNew,
			CreatedAt: time.Now(),
			LastEdit:  time.Now(),
		}

		isValid := order.IsValidStatus(service.OrderStatus(MockOrderStatusPaid))
		assert.True(t, isValid, "Status transition from New to Paid should be valid")
	})

	t.Run("Test Invalid Status Transition from New to Processing", func(t *testing.T) {
		order := &Order{
			OID:       uuid.New(),
			TID:       uuid.New(),
			TPrice:    100.0,
			Status:    OrderStatusNew,
			CreatedAt: time.Now(),
			LastEdit:  time.Now(),
		}

		isValid := order.IsValidStatus(service.OrderStatus(MockOrderStatusProcessing))
		assert.False(t, isValid, "Status transition from New to Processing should be invalid")
	})

	t.Run("Test Valid Status Transition from Paid to Processing", func(t *testing.T) {
		order := &Order{
			OID:       uuid.New(),
			TID:       uuid.New(),
			TPrice:    100.0,
			Status:    OrderStatusPaid,
			CreatedAt: time.Now(),
			LastEdit:  time.Now(),
		}

		isValid := order.IsValidStatus(service.OrderStatus(MockOrderStatusProcessing))
		assert.True(t, isValid, "Status transition from Paid to Processing should be valid")
	})

	t.Run("Test Valid Status Transition from Processing to Done", func(t *testing.T) {
		order := &Order{
			OID:       uuid.New(),
			TID:       uuid.New(),
			TPrice:    100.0,
			Status:    OrderStatusProcessing,
			CreatedAt: time.Now(),
			LastEdit:  time.Now(),
		}

		isValid := order.IsValidStatus(service.OrderStatus(MockOrderStatusDone))
		assert.True(t, isValid, "Status transition from Processing to Done should be valid")
	})

	t.Run("Test Invalid Status Transition with Unknown Status", func(t *testing.T) {
		order := &Order{
			OID:       uuid.New(),
			TID:       uuid.New(),
			TPrice:    100.0,
			Status:    OrderStatusProcessing,
			CreatedAt: time.Now(),
			LastEdit:  time.Now(),
		}

		isValid := order.IsValidStatus(service.OrderStatus(MockOrderStatusInvalid))
		assert.False(t, isValid, "Transition to an unknown status should be invalid")
	})
}
