package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"memrizr/account/entity"
	"memrizr/account/service/repository"
	"testing"
)

func setupStockTest() (*repository.MockStockRepository, StockService) {
	stockRepo := new(repository.MockStockRepository)
	stockService := NewStockService(&StockConfig{
		StockRepository: stockRepo,
	})
	return stockRepo, stockService
}

func TestStockService(t *testing.T) {
	ctx := context.Background()

	t.Run("GetStockByID", func(t *testing.T) {
		mockRepo, stockSvc := setupStockTest()
		stockID := 1
		stock := &entity.Stock{SID: stockID, Quantity: 100}

		// Setup expectation
		mockRepo.On("GetStockByProductID", ctx, stockID).Return(stock, nil)

		t.Run("success", func(t *testing.T) {
			result, err := stockSvc.GetStockByID(ctx, stockID)

			assert.NoError(t, err)
			assert.Equal(t, stock, result)
			mockRepo.AssertExpectations(t)
		})

		t.Run("not found", func(t *testing.T) {
			nonExistentStockID := 2
			mockRepo.On("GetStockByProductID", ctx, nonExistentStockID).Return(nil, errors.New("stock not found"))

			result, err := stockSvc.GetStockByID(ctx, nonExistentStockID)

			assert.Error(t, err)
			assert.Nil(t, result)
			mockRepo.AssertExpectations(t)
		})
	})

	t.Run("UpdateStockById", func(t *testing.T) {
		mockRepo, stockSvc := setupStockTest()

		stock := &entity.Stock{SID: 1, Quantity: 50}
		updatedStock := &entity.Stock{SID: 1, Quantity: 75}

		t.Run("success", func(t *testing.T) {
			// Setup expectation
			mockRepo.On("UpdateStock", ctx, stock).Return(updatedStock, nil)

			result, err := stockSvc.UpdateStockById(ctx, stock)

			assert.NoError(t, err)
			assert.Equal(t, updatedStock, result)
			mockRepo.AssertExpectations(t)
		})

		t.Run("update error", func(t *testing.T) {
			mockRepo.ExpectedCalls = nil

			mockRepo.On("UpdateStock", ctx, stock).Return(nil, errors.New("update failed"))

			result, err := stockSvc.UpdateStockById(ctx, stock)

			assert.Error(t, err)
			assert.Nil(t, result)
			mockRepo.AssertExpectations(t)
		})
	})
}
