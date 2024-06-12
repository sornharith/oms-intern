package service

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"memrizr/account/entity"
	service "memrizr/account/service/model"
	"strconv"
	"strings"
	"time"
)

type createOrderUsecase struct {
	calPriceRepo entity.CalPriceRepository
	orderRepo    entity.OrderRepository
	stockRepo    entity.StockRepository
}

type CreateOrderconfig struct {
	CalPriceRepo entity.CalPriceRepository
	OrderRepo    entity.OrderRepository
	StockRepo    entity.StockRepository
}

var orderId = 3

func NewCreateOrderUsecase(c *CreateOrderconfig) entity.OrderService {
	return &createOrderUsecase{
		calPriceRepo: c.CalPriceRepo,
		orderRepo:    c.OrderRepo,
		stockRepo:    c.StockRepo,
	}
}

func (u *createOrderUsecase) CreateOrder(tID uuid.UUID) (*entity.Order, error) {
	// Fetch transaction details
	calPrice, err := u.calPriceRepo.GetByID(tID)
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

	// for rollback deduction
	type deduction struct {
		ProductID int
		Amount    int
	}
	var deductions []deduction

	for _, item := range userSelect {
		productID := item.ProductID
		amount := item.Amount

		if err := u.stockRepo.DeductStock(productID, amount); err != nil {
			// Rollback deductions
			for _, d := range deductions {
				if err := u.stockRepo.AddStock(d.ProductID, d.Amount); err != nil {
					log.Printf("failed to rollback stock for product ID %d: %s", d.ProductID, err)
				}
			}
			log.Println("insufficient stock for product ID: " + strconv.Itoa(productID))
			return nil, errors.New("please check your stock")
		}
		// Record successful deduction
		deductions = append(deductions, deduction{ProductID: productID, Amount: amount})
	}

	// Create the order
	order := &entity.Order{
		TID:       tID,
		TPrice:    calPrice.TPrice,
		Status:    entity.OrderStatusNew,
		CreatedAt: time.Now(),
		LastEdit:  time.Now(),
	}

	if err := u.orderRepo.CreateOrder(order); err != nil {
		return nil, err
	}

	return order, nil
}
func (u *createOrderUsecase) GetOrderByID(id uuid.UUID) (*entity.Order, error) {
	return u.orderRepo.GetByID(id)
}

func (u *createOrderUsecase) UpdateOrderStatus(o_id uuid.UUID, status string) error {
	order, err := u.orderRepo.GetByID(o_id)
	if err != nil {
		log.Printf("error getting order by id %d", o_id)
		return errors.New("invalid input")
	}
	if order.Status == entity.OrderStatusNew && status == entity.OrderStatusPaid {
		order.Status = status
	} else if isValidStatus(service.OrderStatus(status)) {
		order.Status = status
	} else {
		return errors.New("invalid order status")
	}

	err = u.orderRepo.Update(order)
	if err != nil {
		return err
	}
	return nil
}

func isValidStatus(status service.OrderStatus) bool {
	return status == entity.OrderStatusPaid ||
		status == entity.OrderStatusProcessing ||
		status == entity.OrderStatusDone
}

func (u *createOrderUsecase) DeleteOrder(id int) error {
	//TODO implement me
	panic("implement me")
}
