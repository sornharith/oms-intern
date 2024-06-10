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
	// Get transaction details
	calPrice, err := u.calPriceRepo.GetByID(tID)
	if err != nil {
		return nil, err
	}
	unescaped := strings.Trim(calPrice.UserSelect, "\"")
	user_select := strings.ReplaceAll(unescaped, `\"`, `"`)

	// Check stock for each product in user select
	var userSelect []map[string]interface{}
	if err := json.Unmarshal([]byte(user_select), &userSelect); err != nil {
		log.Printf(`error unmarshalling user select %s`, err)
		return nil, errors.New("unable to parse user select")
	}

	// TODO: check race condition by update the date to check
	for _, item := range userSelect {
		productID := int(item["product_id"].(float64))
		amount := int(item["amount"].(float64))

		stock, err := u.stockRepo.GetStockByProductID(productID)
		if err != nil {
			return nil, err
		}

		if stock.Quantity < amount {
			return nil, errors.New(`insufficient stock for product ID: ` + strconv.Itoa(productID) + ``)
		}
	}

	// Create order
	order := &model.Order{
		TID:       tID,
		TPrice:    float64(int(calPrice.TPrice)),
		Status:    "New",
		CreatedAt: time.Now(),
		LastEdit:  time.Now(),
	}
	if err := u.orderRepo.CreateOrder(order); err != nil {
		return nil, err
	}
	// Deduct stock
	for _, item := range userSelect {
		productID := int(item["product_id"].(float64))
		amount := int(item["amount"].(float64))

		if err := u.stockRepo.DeductStock(productID, amount); err != nil {
			return nil, err
		}
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
