package service

import (
	"context"
	"memrizr/account/entity"
	"memrizr/account/repository"
)

type stockService struct {
	StockRepository repository.StockRepository
}

type StockConfig struct {
	StockRepository repository.StockRepository
}

func NewStockService(c *StockConfig) StockService {
	return &stockService{
		StockRepository: c.StockRepository,
	}
}

func (s stockService) GetStockByID(ctx context.Context, id int) (*entity.Stock, error) {
	return s.StockRepository.GetStockByProductID(ctx, id)
}

func (s stockService) UpdateStockById(ctx context.Context, stock *entity.Stock) (*entity.Stock, error) {
	return s.StockRepository.UpdateStock(ctx, stock)
}
