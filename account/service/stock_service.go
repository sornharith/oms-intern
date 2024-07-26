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
	ctx, span := tracer.Start(ctx, "service get-stock-by-id")
	defer span.End()

	return s.StockRepository.GetStockByProductID(ctx, id)
}

func (s stockService) UpdateStockById(ctx context.Context, stock *entity.Stock) (*entity.Stock, error) {
	ctx, span := tracer.Start(ctx, "service udpate-stock-by-id")
	defer span.End()

	return s.StockRepository.UpdateStock(ctx, stock)
}
