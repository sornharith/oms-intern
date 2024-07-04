package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"memrizr/account/entity"
	"memrizr/account/service/repository"
	"testing"
	"time"
)

func setupOrderTest() (*repository.MockOrderRepository, *repository.MockCalPriceRepository, *repository.MockStockRepository, entity.OrderService) {
	orderRepo := new(repository.MockOrderRepository)
	calPriceRepo := new(repository.MockCalPriceRepository)
	stockRepo := new(repository.MockStockRepository)
	orderservice := NewCreateOrderUsecase(&CreateOrderconfig{
		CalPriceRepo: calPriceRepo,
		OrderRepo:    orderRepo,
		StockRepo:    stockRepo,
	})
	return orderRepo, calPriceRepo, stockRepo, orderservice
}

func TestCreateOrderUsecase(t *testing.T) {
	ctx := context.TODO()
	t.Run("Successful CreateOrder", func(t *testing.T) {
		mockTime := time.Date(2024, time.July, 2, 15, 38, 27, 0, time.UTC)

		mockOrderRepo, mockCalPriceRepo, mockStockRepo, service := setupOrderTest()
		tID := uuid.New()
		calPrice := &entity.CalPrice{
			TID:        tID,
			TPrice:     100.0,
			UserSelect: `[{"product_id": 1, "amount": 2}]`,
		}
		_ = &entity.Order{
			TID:       tID,
			TPrice:    100.0,
			Status:    entity.OrderStatusNew,
			CreatedAt: mockTime,
			LastEdit:  mockTime,
		}
		deductions := map[int]int{1: 2}

		mockCalPriceRepo.On("GetByID", ctx, tID).Return(calPrice, nil)
		mockStockRepo.On("DeductStockBulk", ctx, deductions).Return(nil)
		mockOrderRepo.On("CreateOrder", ctx, mock.AnythingOfType("*entity.Order")).Return(nil)

		result, err := service.CreateOrder(ctx, tID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 100.0, result.TPrice)
		mockCalPriceRepo.AssertExpectations(t)
		mockStockRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Failed to Get CalPrice", func(t *testing.T) {
		_, mockCalPriceRepo, _, service := setupOrderTest()

		tID := uuid.New()
		mockCalPriceRepo.On("GetByID", ctx, tID).Return(nil, errors.New("calPrice not found"))

		result, err := service.CreateOrder(ctx, tID)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockCalPriceRepo.AssertExpectations(t)
	})

	t.Run("Successful GetOrderByID", func(t *testing.T) {
		mockOrderRepo, _, _, service := setupOrderTest()

		tID := uuid.New()
		order := &entity.Order{
			TID:    tID,
			TPrice: 100.0,
			Status: entity.OrderStatusNew,
		}

		mockOrderRepo.On("GetByID", ctx, tID).Return(order, nil)

		result, err := service.GetOrderByID(ctx, tID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 100.0, result.TPrice)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Failed to Get Order", func(t *testing.T) {
		mockOrderRepo, _, _, service := setupOrderTest()

		tID := uuid.New()
		mockOrderRepo.On("GetByID", ctx, tID).Return(nil, errors.New("order not found"))

		result, err := service.GetOrderByID(ctx, tID)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Successful UpdateOrderStatus", func(t *testing.T) {
		mockOrderRepo, _, _, service := setupOrderTest()

		oID := uuid.New()
		order := &entity.Order{
			TID:    oID,
			TPrice: 100.0,
			Status: entity.OrderStatusNew,
		}

		mockOrderRepo.On("GetByID", ctx, oID).Return(order, nil)
		mockOrderRepo.On("Update", ctx, mock.AnythingOfType("*entity.Order")).Return(nil)

		err := service.UpdateOrderStatus(ctx, oID, entity.OrderStatusPaid)

		assert.NoError(t, err)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Failed to Update Order Status", func(t *testing.T) {
		mockOrderRepo, _, _, service := setupOrderTest()
		oID := uuid.New()
		order := &entity.Order{
			TID:    oID,
			TPrice: 100.0,
			Status: entity.OrderStatusNew,
		}

		mockOrderRepo.On("GetByID", ctx, oID).Return(order, nil)
		mockOrderRepo.On("Update", ctx, mock.AnythingOfType("*entity.Order")).Return(errors.New("update failed"))

		err := service.UpdateOrderStatus(ctx, oID, entity.OrderStatusPaid)

		assert.Error(t, err)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Invalid Status Update", func(t *testing.T) {
		mockOrderRepo, _, _, service := setupOrderTest()

		oID := uuid.New()
		order := &entity.Order{
			TID:    oID,
			TPrice: 100.0,
			Status: entity.OrderStatusNew,
		}

		mockOrderRepo.On("GetByID", ctx, oID).Return(order, nil)

		err := service.UpdateOrderStatus(ctx, oID, "InvalidStatus")

		assert.Error(t, err)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Successful DeleteOrder", func(t *testing.T) {
		mockOrderRepo, _, _, service := setupOrderTest()

		orderID := 1
		mockOrderRepo.On("Delete", ctx, orderID).Return(nil)

		err := service.DeleteOrder(ctx, orderID)

		assert.NoError(t, err)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Failed to Delete Order", func(t *testing.T) {
		mockOrderRepo, _, _, service := setupOrderTest()

		orderID := 1
		mockOrderRepo.On("Delete", ctx, orderID).Return(errors.New("delete failed"))

		err := service.DeleteOrder(ctx, orderID)

		assert.Error(t, err)
		mockOrderRepo.AssertExpectations(t)
	})

}
