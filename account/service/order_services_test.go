package service

import (
	"testing"
	"time"
)

type MockTime struct{}

var fixedTime = time.Date(2024, time.July, 2, 15, 5, 39, 0, time.Local)

func (m *MockTime) Now() time.Time {
	return fixedTime
}
func TestCreateOrderUsecase(t *testing.T) {
	//ctx := context.TODO()
	// Replace time.Now() with MockTime
	//timeNow := &MockTime{}

	//t.Run("Test CreateOrder Success", func(t *testing.T) {
	//	mockOrderRepo := new(repository.MockOrderRepository)
	//	mockCalPriceRepo := new(repository.MockCalPriceRepository)
	//	mockStockRepo := new(repository.MockStockRepository)
	//
	//	tID := uuid.New()
	//	calPrice := &entity.CalPrice{
	//		TID:        tID,
	//		TPrice:     100.0,
	//		UserSelect: `[{"product_id": 1, "amount": 2}]`,
	//		Address:    "International",
	//	}
	//	order := &entity.Order{
	//		TID:       tID,
	//		TPrice:    100.0,
	//		Status:    entity.OrderStatusNew,
	//		CreatedAt: time.Now(),
	//		LastEdit:  time.Now(),
	//	}
	//	deductions := map[int]int{1: 2}
	//
	//	mockCalPriceRepo.On("GetByID", ctx, tID).Return(calPrice, nil)
	//	mockStockRepo.On("DeductStockBulk", ctx, deductions).Return(nil)
	//	mockOrderRepo.On("CreateOrder", ctx, order).Return(nil)
	//
	//	service := NewCreateOrderUsecase(&CreateOrderconfig{
	//		CalPriceRepo: mockCalPriceRepo,
	//		OrderRepo:    mockOrderRepo,
	//		StockRepo:    mockStockRepo,
	//	})
	//
	//	// Call the CreateOrder method
	//	result, err := service.CreateOrder(ctx, tID)
	//	t.Log(result)
	//	// Validate the results
	//	assert.NoError(t, err)
	//	assert.NotNil(t, result)
	//	assert.Equal(t, 200.0, result.TPrice)
	//	mockCalPriceRepo.AssertExpectations(t)
	//	mockStockRepo.AssertExpectations(t)
	//	mockOrderRepo.AssertExpectations(t)
	//})

	//t.Run("Test CreateOrder CalPrice GetByID Error", func(t *testing.T) {
	//	mockOrderRepo := new(repository.MockOrderRepository)
	//	mockCalPriceRepo := new(repository.MockCalPriceRepository)
	//	mockStockRepo := new(repository.MockStockRepository)
	//
	//	tID := uuid.New()
	//
	//	mockCalPriceRepo.On("GetByID", ctx, tID).Return(nil, errors.New("not found"))
	//
	//	service := NewCreateOrderUsecase(&CreateOrderconfig{
	//		CalPriceRepo: mockCalPriceRepo,
	//		OrderRepo:    mockOrderRepo,
	//		StockRepo:    mockStockRepo,
	//	})
	//
	//	result, err := service.CreateOrder(ctx, tID)
	//
	//	assert.Error(t, err)
	//	assert.Nil(t, result)
	//	assert.Equal(t, "not found", err.Error())
	//	mockCalPriceRepo.AssertExpectations(t)
	//})
	//
	//t.Run("Test CreateOrder Stock DeductStockBulk Error", func(t *testing.T) {
	//	mockOrderRepo := new(repository.MockOrderRepository)
	//	mockCalPriceRepo := new(repository.MockCalPriceRepository)
	//	mockStockRepo := new(repository.MockStockRepository)
	//
	//	tID := uuid.New()
	//	calPrice := &entity.CalPrice{
	//		TID:        tID,
	//		TPrice:     100.0,
	//		UserSelect: `[{"product_id": 1, "amount": 2}]`,
	//	}
	//	_ = []entity.UserSelectItem{
	//		{ProductID: 1, Amount: 2},
	//	}
	//	deductions := map[int]int{1: 2}
	//
	//	mockCalPriceRepo.On("GetByID", ctx, tID).Return(calPrice, nil)
	//	mockStockRepo.On("DeductStockBulk", ctx, deductions).Return(errors.New("stock error"))
	//
	//	service := NewCreateOrderUsecase(&CreateOrderconfig{
	//		CalPriceRepo: mockCalPriceRepo,
	//		OrderRepo:    mockOrderRepo,
	//		StockRepo:    mockStockRepo,
	//	})
	//
	//	result, err := service.CreateOrder(ctx, tID)
	//
	//	assert.Error(t, err)
	//	assert.Nil(t, result)
	//	assert.Equal(t, "stock error", err.Error())
	//	mockCalPriceRepo.AssertExpectations(t)
	//	mockStockRepo.AssertExpectations(t)
	//})
	//
	//t.Run("Test CreateOrder Order Creation Error", func(t *testing.T) {
	//	mockOrderRepo := new(repository.MockOrderRepository)
	//	mockCalPriceRepo := new(repository.MockCalPriceRepository)
	//	mockStockRepo := new(repository.MockStockRepository)
	//
	//	tID := uuid.New()
	//	calPrice := &entity.CalPrice{
	//		TID:        tID,
	//		TPrice:     100.0,
	//		UserSelect: `[{"product_id": 1, "amount": 2}]`,
	//	}
	//	order := &entity.Order{
	//		TID:       tID,
	//		TPrice:    100.0,
	//		Status:    entity.OrderStatusNew,
	//		CreatedAt: time.Now(),
	//		LastEdit:  time.Now(),
	//	}
	//	_ = []entity.UserSelectItem{
	//		{ProductID: 1, Amount: 2},
	//	}
	//	_ = []map[string]interface{}{
	//		{"product_id": float64(1), "amount": float64(2)},
	//	}
	//	deductions := map[int]int{1: 2}
	//
	//	mockCalPriceRepo.On("GetByID", ctx, tID).Return(calPrice, nil)
	//	mockStockRepo.On("DeductStockBulk", ctx, deductions).Return(nil)
	//	mockOrderRepo.On("CreateOrder", ctx, order).Return(errors.New("order creation error"))
	//
	//	service := NewCreateOrderUsecase(&CreateOrderconfig{
	//		CalPriceRepo: mockCalPriceRepo,
	//		OrderRepo:    mockOrderRepo,
	//		StockRepo:    mockStockRepo,
	//	})
	//
	//	result, err := service.CreateOrder(ctx, tID)
	//
	//	assert.Error(t, err)
	//	assert.Nil(t, result)
	//	assert.Equal(t, "order creation error", err.Error())
	//	mockCalPriceRepo.AssertExpectations(t)
	//	mockStockRepo.AssertExpectations(t)
	//	mockOrderRepo.AssertExpectations(t)
	//})
}
