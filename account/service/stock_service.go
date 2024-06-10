package service

import (
	"memrizr/account/model"
)

type stockService struct {
	StockRepository model.StockRepository
}

type StockConfig struct {
	StockRepository model.StockRepository
}

func NewStockService(c *StockConfig) model.StockService {
	return &stockService{
		StockRepository: c.StockRepository,
	}
}

func (s stockService) GetStockByID(id int) (*model.Stock, error) {
	return s.StockRepository.GetStockByProductID(id)
}
