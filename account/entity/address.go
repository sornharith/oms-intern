package entity

import apperror "memrizr/account/entity/apperrors"

const (
	AddressDomestic      = "Domestic"
	AddressInternational = "International"
	PriceDomestic        = 50.0
	PriceInternational   = 100.0
)

type Address struct {
	Address string
}

func (a *Address) Init(add string) {
	a.Address = add
}

func (a *Address) Price() (float64, error) {
	var Price float64
	switch a.Address {
	case AddressDomestic:
		Price += PriceDomestic
	case AddressInternational:
		Price += PriceInternational
	default:
		return 0, apperror.NewBadRequest("Incorrect address")
	}
	return Price, nil
}
