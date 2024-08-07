package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"memrizr/account/entity"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *MockOrderRepository) CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	args := m.Called(ctx, order)
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	args := m.Called(ctx, order)
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *MockOrderRepository) Delete(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Order), args.Error(1)
}
