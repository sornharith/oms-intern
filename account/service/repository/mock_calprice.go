package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"memrizr/account/entity"
)

type MockCalPriceRepository struct {
	mock.Mock
}

func (m *MockCalPriceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.CalPrice, error) {
	args := m.Called(ctx, id)
	calPrice, _ := args.Get(0).(*entity.CalPrice)
	return calPrice, args.Error(1)
}

func (m *MockCalPriceRepository) Update(ctx context.Context, calPrice *entity.CalPrice) error {
	args := m.Called(ctx, calPrice)
	return args.Error(0)
}

func (m *MockCalPriceRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCalPriceRepository) CalculateTotalPrice(ctx context.Context, userSelect []map[string]interface{}) (float64, error) {
	args := m.Called(ctx, userSelect)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockCalPriceRepository) CreateCalPrice(ctx context.Context, calPrice *entity.CalPrice) error {
	args := m.Called(ctx, calPrice)
	return args.Error(0)
}
