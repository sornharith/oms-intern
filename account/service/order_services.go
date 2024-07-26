package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"memrizr/account/entity"
	"memrizr/account/repository"
	service "memrizr/account/service/model"
	"strings"
)

type createOrderUsecase struct {
	calPriceRepo repository.CalPriceRepository
	orderRepo    repository.OrderRepository
	stockRepo    repository.StockRepository
}

type CreateOrderconfig struct {
	CalPriceRepo repository.CalPriceRepository
	OrderRepo    repository.OrderRepository
	StockRepo    repository.StockRepository
}

func NewCreateOrderUsecase(c *CreateOrderconfig) OrderService {
	return &createOrderUsecase{
		calPriceRepo: c.CalPriceRepo,
		orderRepo:    c.OrderRepo,
		stockRepo:    c.StockRepo,
	}
}

func (u *createOrderUsecase) CreateOrder(ctx context.Context, tID uuid.UUID) (*entity.Order, error) {
	ctx, span := tracer.Start(ctx, "service get-stock-by-id")
	defer span.End()
	// Fetch transaction details
	calPrice, err := u.calPriceRepo.GetByID(ctx, tID)
	if err != nil {
		return nil, err
	}

	// Parse user_select
	unescaped := strings.Trim(calPrice.UserSelect, "\"")
	userSelectJSON := strings.ReplaceAll(unescaped, `\"`, `"`)

	var userSelect []entity.UserSelectItem
	if err := json.Unmarshal([]byte(userSelectJSON), &userSelect); err != nil {
		return nil, errors.New("unable to parse user select: " + err.Error())
	}

	// Prepare the deductions
	deductions := make(map[int]int)
	for _, item := range userSelect {
		productID := item.ProductID
		amount := item.Amount
		deductions[productID] = amount
	}

	// Attempt to deduct stock in bulk
	if err := u.stockRepo.DeductStockBulk(ctx, deductions); err != nil {
		return nil, err
	}
	// Create the order
	order := &entity.Order{
		TID:    tID,
		TPrice: calPrice.TPrice,
		Status: entity.OrderStatusNew,
	}

	res, err := u.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (u *createOrderUsecase) GetOrderByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	ctx, span := tracer.Start(ctx, "service get-stock-by-id")
	defer span.End()

	return u.orderRepo.GetByID(ctx, id)
}

func (u *createOrderUsecase) UpdateOrderStatus(ctx context.Context, o_id uuid.UUID, status string) (*entity.Order, error) {
	ctx, span := tracer.Start(ctx, "service update-order-status")
	defer span.End()

	order, err := u.orderRepo.GetByID(ctx, o_id)
	if err != nil {
		log.Printf("error getting order by id %d", o_id)
		return nil, errors.New("invalid input")
	}

	if order.Status == entity.OrderStatusNew && status == entity.OrderStatusPaid {
		order.Status = status
	} else if order.IsValidStatus(service.OrderStatus(status)) {
		log.Printf("status %s", status)
		order.Status = status
	} else {
		return nil, errors.New("invalid order status")
	}

	response, err := u.orderRepo.Update(ctx, order)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (u *createOrderUsecase) DeleteOrder(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	return u.orderRepo.Delete(ctx, id)
}
