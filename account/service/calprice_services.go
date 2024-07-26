package service

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"memrizr/account/entity"
	"memrizr/account/repository"
)

type calPriceUsecase struct {
	calPriceRepo repository.CalPriceRepository
	productRepo  repository.ProductRepository
}
type CalpConfig struct {
	CalPriceRepo repository.CalPriceRepository
	ProductRepo  repository.ProductRepository
}

func NewCalPriceUsecase(c *CalpConfig) CalPriceService {
	return &calPriceUsecase{
		calPriceRepo: c.CalPriceRepo,
		productRepo:  c.ProductRepo,
	}
}

func (u *calPriceUsecase) GetCalPriceByID(ctx context.Context, id uuid.UUID) (*entity.CalPrice, error) {
	cts, span := tracer.Start(ctx, "service get-stock-by-id")
	defer span.End()
	
	res, err := u.calPriceRepo.GetByID(cts, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *calPriceUsecase) UpdateCalPrice(ctx context.Context, calPrice *entity.CalPrice) (*entity.CalPrice, error) {
	return u.calPriceRepo.Update(ctx, calPrice)
}

func (u *calPriceUsecase) DeleteCalPrice(ctx context.Context, id uuid.UUID) (*entity.CalPrice, error) {
	return u.calPriceRepo.Delete(ctx, id)
}

func (u *calPriceUsecase) CreateCalPrice(ctx context.Context, calPrice *entity.CalPrice) (*entity.CalPrice, error) {
	ctx, span := tracer.Start(ctx, "service create-calprice")
	defer span.End()
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
	res, err := u.calPriceRepo.CreateCalPrice(ctx, calPrice)
	if err != nil {
		return nil, err
	}
	return res, nil
}
