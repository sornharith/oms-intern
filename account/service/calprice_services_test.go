package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"memrizr/account/entity"
	"memrizr/account/service/repository"
	"testing"
)

func TestCalPriceUsecase_Get(t *testing.T) {
	ctx := context.TODO()

	t.Run("Test GetCalPriceByID Success", func(t *testing.T) {
		mockCalPriceRepo := new(repository.MockCalPriceRepository)
		id := uuid.New()
		calPrice := &entity.CalPrice{TID: id}

		mockCalPriceRepo.On("GetByID", ctx, id).Return(calPrice, nil)

		service := NewCalPriceUsecase(&CalpConfig{CalPriceRepo: mockCalPriceRepo})
		result, err := service.GetCalPriceByID(ctx, id)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, calPrice, result)
		mockCalPriceRepo.AssertExpectations(t)
	})

	t.Run("Test GetCalPriceByID Error", func(t *testing.T) {
		mockCalPriceRepo := new(repository.MockCalPriceRepository)
		id := uuid.New()

		mockCalPriceRepo.On("GetByID", ctx, id).Return((*entity.CalPrice)(nil), errors.New("not found"))

		service := NewCalPriceUsecase(&CalpConfig{CalPriceRepo: mockCalPriceRepo})
		result, err := service.GetCalPriceByID(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "not found", err.Error())
		mockCalPriceRepo.AssertExpectations(t)
	})

}
func TestCalPriceUsecase_Update(t *testing.T) {
	ctx := context.TODO()
	t.Run("Test UpdateCalPrice Success", func(t *testing.T) {
		mockCalPriceRepo := new(repository.MockCalPriceRepository)
		calPrice := &entity.CalPrice{TID: uuid.New()}

		mockCalPriceRepo.On("Update", ctx, calPrice).Return(nil)

		service := NewCalPriceUsecase(&CalpConfig{CalPriceRepo: mockCalPriceRepo})
		err := service.UpdateCalPrice(ctx, calPrice)

		assert.NoError(t, err)
		mockCalPriceRepo.AssertExpectations(t)
	})

	t.Run("Test UpdateCalPrice Error", func(t *testing.T) {
		mockCalPriceRepo := new(repository.MockCalPriceRepository)
		calPrice := &entity.CalPrice{TID: uuid.New()}

		mockCalPriceRepo.On("Update", ctx, calPrice).Return(errors.New("update failed"))

		service := NewCalPriceUsecase(&CalpConfig{CalPriceRepo: mockCalPriceRepo})
		err := service.UpdateCalPrice(ctx, calPrice)

		assert.Error(t, err)
		assert.Equal(t, "update failed", err.Error())
		mockCalPriceRepo.AssertExpectations(t)
	})
}
func TestCalPriceUsecase_Delete(t *testing.T) {
	ctx := context.TODO()
	t.Run("Test DeleteCalPrice Success", func(t *testing.T) {
		mockCalPriceRepo := new(repository.MockCalPriceRepository)

		mockCalPriceRepo.On("Delete", ctx, 1).Return(nil)

		service := NewCalPriceUsecase(&CalpConfig{CalPriceRepo: mockCalPriceRepo})
		err := service.DeleteCalPrice(ctx, 1)

		assert.NoError(t, err)
		mockCalPriceRepo.AssertExpectations(t)
	})

	t.Run("Test DeleteCalPrice Error", func(t *testing.T) {
		mockCalPriceRepo := new(repository.MockCalPriceRepository)

		mockCalPriceRepo.On("Delete", ctx, 1).Return(errors.New("delete failed"))

		service := NewCalPriceUsecase(&CalpConfig{CalPriceRepo: mockCalPriceRepo})
		err := service.DeleteCalPrice(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, "delete failed", err.Error())
		mockCalPriceRepo.AssertExpectations(t)
	})
}
func TestCalPriceUsecase_Create(t *testing.T) {
	ctx := context.TODO()
	t.Run("Test CreateCalPrice Success", func(t *testing.T) {
		mockCalPriceRepo := new(repository.MockCalPriceRepository)
		calPrice := &entity.CalPrice{
			TID:        uuid.New(),
			TPrice:     20.0,
			UserSelect: `[{"product_id": 1, "amount": 2}]`,
			Address:    "International",
		}
		userSelect := []map[string]interface{}{
			{"product_id": float64(1), "amount": float64(2)}, // Ensure these are float64
		}

		mockCalPriceRepo.On("CalculateTotalPrice", ctx, userSelect).Return(20.0, nil)
		mockCalPriceRepo.On("CreateCalPrice", ctx, calPrice).Return(nil)

		service := NewCalPriceUsecase(&CalpConfig{CalPriceRepo: mockCalPriceRepo})

		result, err := service.CreateCalPrice(ctx, calPrice)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 120.0, result.TPrice)
		mockCalPriceRepo.AssertExpectations(t)
	})
	t.Run("Test CreateCalPrice JSON Error", func(t *testing.T) {
		mockCalPriceRepo := new(repository.MockCalPriceRepository)
		calPrice := &entity.CalPrice{
			TID:        uuid.New(),
			UserSelect: `invalid json`,
		}

		service := NewCalPriceUsecase(&CalpConfig{CalPriceRepo: mockCalPriceRepo})

		result, err := service.CreateCalPrice(ctx, calPrice)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockCalPriceRepo.AssertExpectations(t)
	})
	t.Run("Test CreateCalPrice Calculation Error", func(t *testing.T) {
		mockCalPriceRepo := new(repository.MockCalPriceRepository)
		calPrice := &entity.CalPrice{
			TID:        uuid.New(),
			UserSelect: `[{"product_id": 1, "amount": 2}]`,
		}
		userSelect := []map[string]interface{}{
			{"product_id": float64(1), "amount": float64(2)}, // Ensure these are float64
		}

		mockCalPriceRepo.On("CalculateTotalPrice", ctx, userSelect).Return(0.0, errors.New("calculation failed"))

		service := NewCalPriceUsecase(&CalpConfig{CalPriceRepo: mockCalPriceRepo})

		result, err := service.CreateCalPrice(ctx, calPrice)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "calculation failed", err.Error())
		mockCalPriceRepo.AssertExpectations(t)
	})
	t.Run("Test CreateCalPrice Create Error", func(t *testing.T) {
		mockCalPriceRepo := new(repository.MockCalPriceRepository)
		calPrice := &entity.CalPrice{
			TID:        uuid.New(),
			UserSelect: `[{"product_id": 1, "amount": 2}]`,
			Address:    "International",
		}
		userSelect := []map[string]interface{}{
			{"product_id": float64(1), "amount": float64(2)}, // Ensure these are float64
		}

		mockCalPriceRepo.On("CalculateTotalPrice", ctx, userSelect).Return(20.0, nil)
		mockCalPriceRepo.On("CreateCalPrice", ctx, calPrice).Return(errors.New("create failed"))

		service := NewCalPriceUsecase(&CalpConfig{CalPriceRepo: mockCalPriceRepo})

		result, err := service.CreateCalPrice(ctx, calPrice)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "create failed", err.Error())
		mockCalPriceRepo.AssertExpectations(t)
	})
}
