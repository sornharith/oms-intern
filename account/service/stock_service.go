package service

import (
	"context"
	"memrizr/account/entity"
)

type stockService struct {
	StockRepository entity.StockRepository
}

type StockConfig struct {
	StockRepository entity.StockRepository
}

func NewStockService(c *StockConfig) entity.StockService {
	return &stockService{
		StockRepository: c.StockRepository,
	}
}

func (s stockService) GetStockByID(ctx context.Context, id int) (*entity.Stock, error) {
	return s.StockRepository.GetStockByProductID(ctx, id)
}

func (s stockService) UpdateStockById(ctx context.Context, stock *entity.Stock) error {
	return s.StockRepository.UpdateStock(ctx, stock.SID, stock.Quantity)
}
