package repository

import (
	"context"
	apperror "memrizr/account/entity/apperrors"
	"memrizr/account/observability/logger"
	service "memrizr/account/service/model"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type productRepository struct {
	DB *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepository{
		DB: db,
	}
}

func (p productRepository) GetallProductStock(ctx context.Context) ([]service.ProductStock, error) {
	_, span := tracer.Start(ctx, "repository get-product-with-stock")
	defer span.End()
	var product []service.ProductStock

	query := `SELECT p.p_id AS ProductId, s.s_id AS StockId, s.quantity AS Quantity
              FROM products p
              JOIN stocks s ON p.s_id = s.s_id`
	err := p.DB.Select(&product, query)
	if err != nil {
		logger.LogError(apperror.CusNotFound("not found product", "3044"), "not found product", logrus.Fields{
			"at": "repository",
		})
	}	
	return product, err
}
