package service

import (
	"context"
	"memrizr/account/entity"
	apperror "memrizr/account/entity/apperrors"
	"memrizr/account/observability/logger"
	"memrizr/account/repository"
	"strconv"

	"github.com/sirupsen/logrus"
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

	res, err := s.StockRepository.GetStockByProductID(ctx, id)
	if err != nil {
		logger.LogError(apperror.CusNotFound(strconv.Itoa(id), "2044"), "error from respository", logrus.Fields{
			"at": "service",
		})
	}
	return res, err
}

func (s stockService) UpdateStockById(ctx context.Context, stock *entity.Stock) (*entity.Stock, error) {
	ctx, span := tracer.Start(ctx, "service udpate-stock-by-id")
	defer span.End()

	
	res, err := s.StockRepository.UpdateStock(ctx, stock)
	if err != nil {
		logger.LogError(apperror.CusNotFound(strconv.Itoa(stock.SID), "2044"), "error from respository", logrus.Fields{
			"at": "service",
		})
	}
	return res , err
}
