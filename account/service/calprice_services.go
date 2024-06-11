package service

import (
	"encoding/json"
	"github.com/google/uuid"
	"memrizr/account/model"
	apperror "memrizr/account/model/apperrors"
)

type calPriceUsecase struct {
	calPriceRepo model.CalPriceRepository
	productRepo  model.ProductRepository
}
type CalpConfig struct {
	CalPriceRepo model.CalPriceRepository
	ProductRepo  model.ProductRepository
}

func NewCalPriceUsecase(c *CalpConfig) model.CalPriceService {
	return &calPriceUsecase{
		calPriceRepo: c.CalPriceRepo,
		productRepo:  c.ProductRepo,
	}
}

func (u *calPriceUsecase) GetCalPriceByID(id uuid.UUID) (*model.CalPrice, error) {
	return u.calPriceRepo.GetByID(id)
}

func (u *calPriceUsecase) UpdateCalPrice(calPrice *model.CalPrice) error {
	return u.calPriceRepo.Update(calPrice)
}

func (u *calPriceUsecase) DeleteCalPrice(id int) error {
	return u.calPriceRepo.Delete(id)
}

func (u *calPriceUsecase) CreateCalPrice(calPrice *model.CalPrice) (*model.CalPrice, error) {
	// Convert UserSelect JSON string back to []map[string]interface{}
	var userSelect []map[string]interface{}
	if err := json.Unmarshal([]byte(calPrice.UserSelect), &userSelect); err != nil {
		return nil, err
	}

	// Calculate total price
	totalPrice, err := u.calPriceRepo.CalculateTotalPrice(userSelect)
	if err != nil {
		return nil, err
	}
	// adding the address price to the total price
	switch calPrice.Address {
	case model.AddressDomestic:
		totalPrice += model.PriceDomestic
	case model.AddressInternational:
		totalPrice += model.PriceInternational
	default:
		return nil, apperror.NewBadRequest("Incorrect address")
	}

	// Update the CalPrice entity with the calculated total price
	calPrice.TPrice = totalPrice
	// Store in repository
	if err := u.calPriceRepo.CreateCalPrice(calPrice); err != nil {
		return nil, err
	}
	return calPrice, nil
}
