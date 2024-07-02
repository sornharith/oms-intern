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
	return args.Get(0).(*entity.Stock), args.Error(1)
}

func (m *MockStockRepository) DeductStockBulk(ctx context.Context, deductions map[int]int) error {
	args := m.Called(ctx, deductions)
	return args.Error(0)
}

func (m *MockStockRepository) AddStock(ctx context.Context, productID int, amount int) error {
	args := m.Called(ctx, productID, amount)
	return args.Error(0)
}

func (m *MockStockRepository) UpdateStock(ctx context.Context, stockID int, amount int) error {
	args := m.Called(ctx, stockID, amount)
	return args.Error(0)
}
