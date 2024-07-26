package service

import (
	"context"
	"memrizr/account/entity"
	service "memrizr/account/service/model"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type CalPriceService interface {
	GetCalPriceByID(ctx context.Context, id uuid.UUID) (*entity.CalPrice, error)
	CreateCalPrice(ctx context.Context, userSelect *entity.CalPrice) (*entity.CalPrice, error)
	UpdateCalPrice(ctx context.Context, calPrice *entity.CalPrice) (*entity.CalPrice, error)
	DeleteCalPrice(ctx context.Context, id uuid.UUID) (*entity.CalPrice, error)
}

type OrderService interface {
	GetOrderByID(ctx context.Context, id uuid.UUID) (*entity.Order, error)
	CreateOrder(ctx context.Context, tID uuid.UUID) (*entity.Order, error)
	UpdateOrderStatus(ctx context.Context, o_id uuid.UUID, status string) (*entity.Order, error)
	DeleteOrder(ctx context.Context, id uuid.UUID) (*entity.Order, error)
}

type ProductService interface {
	GetallProductwithstock(ctx context.Context) ([]service.ProductStock, error)
}

type StockService interface {
	GetStockByID(ctx context.Context, id int) (*entity.Stock, error)
	UpdateStockById(ctx context.Context, stock *entity.Stock) (*entity.Stock, error)
}

var tracer = otel.GetTracerProvider().Tracer("service_layer")
