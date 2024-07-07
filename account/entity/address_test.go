package entity

import (
	apperrors "memrizr/account/entity/apperrors"
	"testing"
)

func TestAddress_Price(t *testing.T) {
	tests := []struct {
		name          string
		address       string
		expectedPrice float64
		expectedError error
	}{
		{
			name:          "Domestic Address",
			address:       AddressDomestic,
			expectedPrice: PriceDomestic,
			expectedError: nil,
		},
		{
			name:          "International Address",
			address:       AddressInternational,
			expectedPrice: PriceInternational,
			expectedError: nil,
		},
		{
			name:          "Invalid Address",
			address:       "Unknown",
			expectedPrice: 0,
			expectedError: apperrors.NewBadRequest("Incorrect address"),
		},
		{
			name:          "Empty Address",
			address:       "",
			expectedPrice: 0,
			expectedError: apperrors.NewBadRequest("Incorrect address"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			addr := &Address{}
			addr.Init(test.address)

			price, err := addr.Price()

			if price != test.expectedPrice {
				t.Errorf("Price() = %v; want %v", price, test.expectedPrice)
			}

			if err != nil && test.expectedError == nil {
				t.Errorf("Price() unexpected error: %v", err)
			} else if err == nil && test.expectedError != nil {
				t.Errorf("Price() expected error: %v, got nil", test.expectedError)
			} else if err != nil && test.expectedError != nil && err.Error() != test.expectedError.Error() {
				t.Errorf("Price() error = %v; want %v", err, test.expectedError)
			}
		})
	}
}
