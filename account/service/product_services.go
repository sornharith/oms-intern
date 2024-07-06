package service

import (
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

func (p productService) GetallProductwithstock() ([]service.ProductStock, error) {
	return p.ProductRepository.GetallProductStock()
}
