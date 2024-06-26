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
	GetCalPriceByID(id uuid.UUID) (*CalPrice, error)
	CreateCalPrice(userSelect *CalPrice) (*CalPrice, error)
	UpdateCalPrice(calPrice *CalPrice) error
	DeleteCalPrice(id int) error
}

type CalPriceRepository interface {
	GetByID(id uuid.UUID) (*CalPrice, error)
	Update(calPrice *CalPrice) error
	Delete(id int) error
	CalculateTotalPrice(userSelect []map[string]interface{}) (float64, error)
	CreateCalPrice(calPrice *CalPrice) error
}

type OrderService interface {
	GetOrderByID(id uuid.UUID) (*Order, error)
	CreateOrder(tID uuid.UUID) (*Order, error)
	UpdateOrderStatus(o_id uuid.UUID, status string) error
	DeleteOrder(id int) error
}

type OrderRepository interface {
	GetByID(id uuid.UUID) (*Order, error)
	Create(order *Order) error
	Update(order *Order) error
	Delete(id int) error
	CreateOrder(order *Order) error
}
type ProductService interface {
	GetallProductwithstock() ([]service.ProductStock, error)
}
type ProductRepository interface {
	GetallProductStock() ([]service.ProductStock, error)
}

type StockService interface {
	GetStockByID(id int) (*Stock, error)
	UpdateStockById(stock *Stock) error
}

type StockRepository interface {
	GetStockByProductID(productID int) (*Stock, error)
	DeductStockBulk(deductions map[int]int) error
	AddStock(productID int, amount int) error
	UpdateStock(StockId int, amount int) error
}
