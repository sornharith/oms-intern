package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"memrizr/account/entity"
	"memrizr/account/service/repository"
	"testing"
	"time"
)

func setupOrderTest() (*repository.MockOrderRepository, *repository.MockCalPriceRepository, *repository.MockStockRepository, OrderService) {
	orderRepo := new(repository.MockOrderRepository)
	calPriceRepo := new(repository.MockCalPriceRepository)
	stockRepo := new(repository.MockStockRepository)

	orderService := NewCreateOrderUsecase(&CreateOrderconfig{
		CalPriceRepo: calPriceRepo,
		OrderRepo:    orderRepo,
		StockRepo:    stockRepo,
	})

	return orderRepo, calPriceRepo, stockRepo, orderService
}

func TestCreateOrder(t *testing.T) {
	ctx := context.TODO()
	t.Run("success", func(t *testing.T) {

		mockOrderRepo, mockCalPriceRepo, mockStockRepo, usecase := setupOrderTest()
		tID := uuid.New()
		calPrice := &entity.CalPrice{TPrice: 100, UserSelect: `[{"ProductID": 1, "Amount": 2}]`}
		order := &entity.Order{OID: uuid.New(), TID: tID, TPrice: 100, Status: "new", CreatedAt: time.Now(), LastEdit: time.Now()}

		// Mock expectations
		mockCalPriceRepo.On("GetByID", ctx, tID).Return(calPrice, nil)
		mockStockRepo.On("DeductStockBulk", ctx, mock.Anything).Return(nil)
		mockOrderRepo.On("CreateOrder", ctx, mock.Anything).Return(order, nil)

		// Execution
		result, err := usecase.CreateOrder(ctx, tID)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, order, result)

		// Verify
		mockCalPriceRepo.AssertExpectations(t)
		mockStockRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("failure - cal price repo error", func(t *testing.T) {
		// Setup
		_, mockCalPriceRepo, _, usecase := setupOrderTest()

		tID := uuid.New()

		// Mock expectations
		mockCalPriceRepo.On("GetByID", ctx, tID).Return((*entity.CalPrice)(nil), errors.New("cal price error"))

		// Execution
		result, err := usecase.CreateOrder(ctx, tID)

		// Assertions
		require.Error(t, err)
		assert.Nil(t, result)

		// Verify
		mockCalPriceRepo.AssertExpectations(t)
	})

	t.Run("TestGetOrderByID success", func(t *testing.T) {
		mockOrderRepo, _, _, usecase := setupOrderTest()

		id := uuid.New()
		order := &entity.Order{OID: id}

		mockOrderRepo.On("GetByID", ctx, id).Return(order, nil)

		result, err := usecase.GetOrderByID(ctx, id)

		require.NoError(t, err)
		assert.Equal(t, order, result)

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("TestGetOrderByID failure - repository error", func(t *testing.T) {
		mockOrderRepo, _, _, usecase := setupOrderTest()
		id := uuid.New()

		mockOrderRepo.On("GetByID", ctx, id).Return((*entity.Order)(nil), errors.New("repository error"))

		result, err := usecase.GetOrderByID(ctx, id)

		require.Error(t, err)
		assert.Nil(t, result)

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("TestUpdateOrderStatus success", func(t *testing.T) {
		mockOrderRepo, _, _, usecase := setupOrderTest()

		id := uuid.New()
		order := &entity.Order{OID: id, Status: entity.OrderStatusNew}

		mockOrderRepo.On("GetByID", ctx, id).Return(order, nil)
		mockOrderRepo.On("Update", ctx, order).Return(order, nil)

		result, err := usecase.UpdateOrderStatus(ctx, id, entity.OrderStatusPaid)

		require.NoError(t, err)
		assert.Equal(t, entity.OrderStatusPaid, result.Status)

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("TestUpdateOrderStatus failure - invalid status", func(t *testing.T) {
		mockOrderRepo, _, _, usecase := setupOrderTest()

		id := uuid.New()
		order := &entity.Order{OID: id, Status: entity.OrderStatusNew}

		mockOrderRepo.On("GetByID", ctx, id).Return(order, nil)

		result, err := usecase.UpdateOrderStatus(ctx, id, "invalid_status")

		require.Error(t, err)
		assert.Nil(t, result)

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("TestDeleteOrder success", func(t *testing.T) {
		mockOrderRepo, _, _, usecase := setupOrderTest()

		id := uuid.New()
		order := &entity.Order{OID: id}

		mockOrderRepo.On("Delete", ctx, id).Return(order, nil)

		result, err := usecase.DeleteOrder(ctx, id)

		require.NoError(t, err)
		assert.Equal(t, order, result)

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("TestDeleteOrder failure - repository error", func(t *testing.T) {
		mockOrderRepo, _, _, usecase := setupOrderTest()

		id := uuid.New()
		mockOrderRepo.On("Delete", ctx, id).Return((*entity.Order)(nil), errors.New("repository error"))

		result, err := usecase.DeleteOrder(ctx, id)

		require.Error(t, err)
		assert.Nil(t, result)

		mockOrderRepo.AssertExpectations(t)
	})
}
