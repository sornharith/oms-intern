package service

import (
	"context"
	apperror "memrizr/account/entity/apperrors"
	"memrizr/account/observability/logger"
	"memrizr/account/repository"
	service "memrizr/account/service/model"

	"github.com/sirupsen/logrus"
)

type productService struct {
	ProductRepository repository.ProductRepository
}

type ProductConfig struct {
	ProductRepository repository.ProductRepository
}

func NewProductService(c *ProductConfig) ProductService {
	return &productService{
		ProductRepository: c.ProductRepository,
	}
}

func (p productService) GetallProductwithstock(ctx context.Context) ([]service.ProductStock, error) {
	ctx, span := tracer.Start(ctx, "service get-product-with-stock")
	defer span.End()

	res, err := p.ProductRepository.GetallProductStock(ctx)
	if err != nil {
		logger.LogError(apperror.CusNotFound("product is empty", "2044"), "error from respository", logrus.Fields{
			"at": "service",
		})
	}
	return res, err
}
