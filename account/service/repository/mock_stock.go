package repository

import (
	"context"
	"github.com/stretchr/testify/mock"
	"memrizr/account/entity"
)

type MockStockRepository struct {
	mock.Mock
}

func (m *MockStockRepository) GetStockByProductID(ctx context.Context, productID int) (*entity.Stock, error) {
	args := m.Called(ctx, productID)
	if stock, ok := args.Get(0).(*entity.Stock); ok {
		return stock, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStockRepository) DeductStockBulk(ctx context.Context, deductions map[int]int) error {
	args := m.Called(ctx, deductions)
	return args.Error(0)
}

func (m *MockStockRepository) UpdateStock(ctx context.Context, stock *entity.Stock) (*entity.Stock, error) {
	args := m.Called(ctx, stock)
	if updatedStock, ok := args.Get(0).(*entity.Stock); ok {
		return updatedStock, args.Error(1)
	}
	return nil, args.Error(1)
}
