package repository

import (
	"github.com/jmoiron/sqlx"
	"memrizr/account/entity"
	service "memrizr/account/service/model"
)

type productRepository struct {
	DB *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) entity.ProductRepository {
	return &productRepository{
		DB: db,
	}
}

func (p productRepository) GetallProductStock() ([]service.ProductStock, error) {

	var product []service.ProductStock

	query := `SELECT p.p_id AS ProductId, s.s_id AS StockId, s.quantity AS Quantity
              FROM products p
              JOIN stocks s ON p.s_id = s.s_id`
	err := p.DB.Select(&product, query)
	return product, err
}
