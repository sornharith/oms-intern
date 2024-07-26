package service

import (
	"context"
	"memrizr/account/repository"
	service "memrizr/account/service/model"
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

	return p.ProductRepository.GetallProductStock(ctx)
}
