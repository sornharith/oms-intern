package repository

import (
	"context"
	"github.com/google/uuid"
	"memrizr/account/entity"
	service "memrizr/account/service/model"
)

type CalPriceRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.CalPrice, error)
	Update(ctx context.Context, calPrice *entity.CalPrice) (*entity.CalPrice, error)
	Delete(ctx context.Context, id uuid.UUID) (*entity.CalPrice, error)
	CalculateTotalPrice(ctx context.Context, userSelect []map[string]interface{}) (float64, error)
	CreateCalPrice(ctx context.Context, calPrice *entity.CalPrice) (*entity.CalPrice, error)
}

type OrderRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error)
	//Create(order *Order) error
	Update(ctx context.Context, order *entity.Order) (*entity.Order, error)
	Delete(ctx context.Context, id uuid.UUID) (*entity.Order, error)
	CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
}

type ProductRepository interface {
	GetallProductStock() ([]service.ProductStock, error)
}

type StockRepository interface {
	GetStockByProductID(ctx context.Context, productID int) (*entity.Stock, error)
	DeductStockBulk(ctx context.Context, deductions map[int]int) error
	//AddStock(ctx context.Context, productID int, amount int) error
	UpdateStock(ctx context.Context, stock *entity.Stock) (*entity.Stock, error)
}
