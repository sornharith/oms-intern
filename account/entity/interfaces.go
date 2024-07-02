package entity

import (
	"context"
	service "memrizr/account/service/model"

	"github.com/google/uuid"
)

// UserService defines methods the handler layer expects
// any service it interacts with to implement
type UserService interface {
	Get(ctx context.Context, uid uuid.UUID) (*User, error)
	Signup(ctx context.Context, u *User) error
}

// UserRepository defines methods the service layer expects
// any repository it interacts with to implement
type UserRepository interface {
	FindByID(ctx context.Context, uid uuid.UUID) (*User, error)
	Create(ctx context.Context, u *User) error
}

type CalPriceService interface {
	GetCalPriceByID(ctx context.Context, id uuid.UUID) (*CalPrice, error)
	CreateCalPrice(ctx context.Context, userSelect *CalPrice) (*CalPrice, error)
	UpdateCalPrice(ctx context.Context, calPrice *CalPrice) error
	DeleteCalPrice(ctx context.Context, id int) error
}

type CalPriceRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*CalPrice, error)
	Update(ctx context.Context, calPrice *CalPrice) error
	Delete(ctx context.Context, id int) error
	CalculateTotalPrice(ctx context.Context, userSelect []map[string]interface{}) (float64, error)
	CreateCalPrice(ctx context.Context, calPrice *CalPrice) error
}

type OrderService interface {
	GetOrderByID(ctx context.Context, id uuid.UUID) (*Order, error)
	CreateOrder(ctx context.Context, tID uuid.UUID) (*Order, error)
	UpdateOrderStatus(ctx context.Context, o_id uuid.UUID, status string) error
	DeleteOrder(ctx context.Context, id int) error
}

type OrderRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Order, error)
	//Create(order *Order) error
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id int) error
	CreateOrder(ctx context.Context, order *Order) error
}
type ProductService interface {
	GetallProductwithstock() ([]service.ProductStock, error)
}
type ProductRepository interface {
	GetallProductStock() ([]service.ProductStock, error)
}

type StockService interface {
	GetStockByID(ctx context.Context, id int) (*Stock, error)
	UpdateStockById(ctx context.Context, stock *Stock) error
}

type StockRepository interface {
	GetStockByProductID(ctx context.Context, productID int) (*Stock, error)
	DeductStockBulk(ctx context.Context, deductions map[int]int) error
	AddStock(ctx context.Context, productID int, amount int) error
	UpdateStock(ctx context.Context, StockId int, amount int) error
}
