package service

import (
	"memrizr/account/entity"
	service "memrizr/account/service/model"
)

type ProductService struct {
	ProductRepository entity.ProductRepository
}

type ProductConfig struct {
	ProductRepository entity.ProductRepository
}

func NewProductService(c *ProductConfig) entity.ProductService {
	return &ProductService{
		ProductRepository: c.ProductRepository,
	}
}

func (p ProductService) GetallProductwithstock() ([]service.ProductStock, error) {
	return p.ProductRepository.GetallProductStock()
}
