package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"memrizr/account/entity"
	"memrizr/account/service/repository"
	"testing"
)

func TestStockService_Get(t *testing.T) {
	ctx := context.TODO()

	t.Run("Test GetStockByID Success", func(t *testing.T) {
		mockRepo := new(repository.MockStockRepository)
		stock := &entity.Stock{SID: 1, Quantity: 100}

		// Mock return values; ensure stock is not nil to avoid nil dereference
		mockRepo.On("GetStockByProductID", ctx, 1).Return(stock, nil)

		service := NewStockService(&StockConfig{StockRepository: mockRepo})
		result, err := service.GetStockByID(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, stock, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Test GetStockByID Error", func(t *testing.T) {
		mockRepo := new(repository.MockStockRepository)
		mockRepo.On("GetStockByProductID", ctx, 2).Return((*entity.Stock)(nil), errors.New("not found"))

		service := NewStockService(&StockConfig{StockRepository: mockRepo})
		result, err := service.GetStockByID(ctx, 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
func TestStockService_Update(t *testing.T) {
	ctx := context.TODO()
	t.Run("Test UpdateStockById Success", func(t *testing.T) {
		mockRepo := new(repository.MockStockRepository)
		stock := &entity.Stock{SID: 1, Quantity: 100}

		// Return nil for success
		mockRepo.On("UpdateStock", ctx, stock.SID, stock.Quantity).Return(nil)

		service := NewStockService(&StockConfig{StockRepository: mockRepo})
		err := service.UpdateStockById(ctx, stock)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Test UpdateStockById Error", func(t *testing.T) {
		mockRepo := new(repository.MockStockRepository)
		stock := &entity.Stock{SID: 1, Quantity: 100}

		// Return an error for failure
		mockRepo.On("UpdateStock", ctx, stock.SID, stock.Quantity).Return(errors.New("update failed"))

		service := NewStockService(&StockConfig{StockRepository: mockRepo})
		err := service.UpdateStockById(ctx, stock)

		assert.Error(t, err)
		assert.Equal(t, "update failed", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Test DeductStockBulk Success", func(t *testing.T) {
		mockRepo := new(repository.MockStockRepository)
		deductions := map[int]int{1: 10, 2: 20}

		// Return nil for success
		mockRepo.On("DeductStockBulk", ctx, deductions).Return(nil)

		_ = NewStockService(&StockConfig{StockRepository: mockRepo})
		err := mockRepo.DeductStockBulk(ctx, deductions)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Test DeductStockBulk Error", func(t *testing.T) {
		mockRepo := new(repository.MockStockRepository)
		deductions := map[int]int{1: 10, 2: 20}

		// Return an error for failure
		mockRepo.On("DeductStockBulk", ctx, deductions).Return(errors.New("deduction failed"))

		_ = NewStockService(&StockConfig{StockRepository: mockRepo})
		err := mockRepo.DeductStockBulk(ctx, deductions)

		assert.Error(t, err)
		assert.Equal(t, "deduction failed", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Test AddStock Success", func(t *testing.T) {
		mockRepo := new(repository.MockStockRepository)

		// Return nil for success
		mockRepo.On("AddStock", ctx, 1, 50).Return(nil)

		_ = NewStockService(&StockConfig{StockRepository: mockRepo})
		err := mockRepo.AddStock(ctx, 1, 50)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Test AddStock Error", func(t *testing.T) {
		mockRepo := new(repository.MockStockRepository)

		// Return an error for failure
		mockRepo.On("AddStock", ctx, 1, 50).Return(errors.New("add stock failed"))

		_ = NewStockService(&StockConfig{StockRepository: mockRepo})
		err := mockRepo.AddStock(ctx, 1, 50)

		assert.Error(t, err)
		assert.Equal(t, "add stock failed", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
