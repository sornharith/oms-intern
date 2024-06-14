package service

import (
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

func (s stockService) GetStockByID(id int) (*entity.Stock, error) {
	return s.StockRepository.GetStockByProductID(id)
}

func (s stockService) UpdateStockById(stock *entity.Stock) error {
	return s.StockRepository.UpdateStock(stock.SID, stock.Quantity)
}
