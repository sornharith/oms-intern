package service

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"

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
	log.Printf("unmarshal userSelect %d", userSelect)

	// Calculate total price
	totalPrice, err := u.calPriceRepo.CalculateTotalPrice(userSelect)
	if err != nil {
		return nil, err
	}

	// TODO: fix the error code from now is 500 to be 400
	if calPrice.Address == "Dominstic" {
		totalPrice += 50
	} else if calPrice.Address == "International" {
		totalPrice += 100
	} else {
		return nil, apperror.NewBadRequest("incorrect address")
	}

	// Update the CalPrice entity with the calculated total price
	calPrice.TPrice = totalPrice
	// Store in repository
	if err := u.calPriceRepo.CreateCalPrice(calPrice); err != nil {
		return nil, err
	}
	return calPrice, nil
}
