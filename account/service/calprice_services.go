package service

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"memrizr/account/entity"
)

type calPriceUsecase struct {
	calPriceRepo entity.CalPriceRepository
	productRepo  entity.ProductRepository
}
type CalpConfig struct {
	CalPriceRepo entity.CalPriceRepository
	ProductRepo  entity.ProductRepository
}

func NewCalPriceUsecase(c *CalpConfig) entity.CalPriceService {
	return &calPriceUsecase{
		calPriceRepo: c.CalPriceRepo,
		productRepo:  c.ProductRepo,
	}
}

func (u *calPriceUsecase) GetCalPriceByID(ctx context.Context, id uuid.UUID) (*entity.CalPrice, error) {
	res, err := u.calPriceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *calPriceUsecase) UpdateCalPrice(ctx context.Context, calPrice *entity.CalPrice) error {
	return u.calPriceRepo.Update(ctx, calPrice)
}

func (u *calPriceUsecase) DeleteCalPrice(ctx context.Context, id int) error {
	return u.calPriceRepo.Delete(ctx, id)
}

func (u *calPriceUsecase) CreateCalPrice(ctx context.Context, calPrice *entity.CalPrice) (*entity.CalPrice, error) {
	// Convert UserSelect JSON string back to []map[string]interface{}
	var userSelect []map[string]interface{}
	if err := json.Unmarshal([]byte(calPrice.UserSelect), &userSelect); err != nil {
		return nil, err
	}

	// Calculate total price
	totalPrice, err := u.calPriceRepo.CalculateTotalPrice(ctx, userSelect)
	if err != nil {
		return nil, err
	}
	// adding the address price to the total price
	Address := new(entity.Address)
	Address.Init(calPrice.Address)

	addprice, err := Address.Price()
	if err != nil {
		return nil, err
	}
	totalPrice += addprice

	// Update the CalPrice entity with the calculated total price
	calPrice.TPrice = totalPrice
	// Store in repository
	if err := u.calPriceRepo.CreateCalPrice(ctx, calPrice); err != nil {
		return nil, err
	}
	return calPrice, nil
}
