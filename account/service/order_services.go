package service

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"memrizr/account/model"
	"strconv"
	"strings"
	"time"
)

type createOrderUsecase struct {
	calPriceRepo model.CalPriceRepository
	orderRepo    model.OrderRepository
	stockRepo    model.StockRepository
}

type CreateOrderconfig struct {
	CalPriceRepo model.CalPriceRepository
	OrderRepo    model.OrderRepository
	StockRepo    model.StockRepository
}

var orderId = 3

func NewCreateOrderUsecase(c *CreateOrderconfig) model.OrderService {
	return &createOrderUsecase{
		calPriceRepo: c.CalPriceRepo,
		orderRepo:    c.OrderRepo,
		stockRepo:    c.StockRepo,
	}
}

func (u *createOrderUsecase) CreateOrder(tID uuid.UUID) (*model.Order, error) {
	// Fetch transaction details
	calPrice, err := u.calPriceRepo.GetByID(tID)
	if err != nil {
		return nil, err
	}

	// Parse user_select
	unescaped := strings.Trim(calPrice.UserSelect, "\"")
	userSelectJSON := strings.ReplaceAll(unescaped, `\"`, `"`)

	var userSelect []model.UserSelectItem
	if err := json.Unmarshal([]byte(userSelectJSON), &userSelect); err != nil {
		log.Printf("error unmarshalling user select %s", err)
		return nil, errors.New("unable to parse user select")
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
			for _, d := range deductions {
				if err := u.stockRepo.AddStock(d.ProductID, d.Amount); err != nil {
					log.Printf("failed to rollback stock for product ID %d: %s", d.ProductID, err)
				}
			}
			return nil, errors.New("insufficient stock for product ID: " + strconv.Itoa(productID))
		}
		// Record successful deduction
		deductions = append(deductions, deduction{ProductID: productID, Amount: amount})
	}

	// Create the order
	order := &model.Order{
		TID:       tID,
		TPrice:    calPrice.TPrice,
		Status:    model.OrderStatusNew,
		CreatedAt: time.Now(),
		LastEdit:  time.Now(),
	}

	if err := u.orderRepo.CreateOrder(order); err != nil {
		return nil, err
	}

	return order, nil
}
func (u *createOrderUsecase) GetOrderByID(id uuid.UUID) (*model.Order, error) {
	return u.orderRepo.GetByID(id)
}

func (u *createOrderUsecase) UpdateOrderStatus(o_id uuid.UUID, status string) error {
	order, err := u.orderRepo.GetByID(o_id)
	if err != nil {
		log.Printf("error getting order by id %d", o_id)
		return errors.New("invalid input")
	}
	if order.Status == "New" && status == "Paid" {
		order.Status = status
	} else if order.Status == "Paid" || order.Status == "Processing" {
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

func (u *createOrderUsecase) DeleteOrder(id int) error {
	//TODO implement me
	panic("implement me")
}
