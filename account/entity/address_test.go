package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddress(t *testing.T) {
	t.Run("Test Init and Price for Domestic Address", func(t *testing.T) {
		addr := &Address{}
		addr.Init(AddressDomestic)

		price, err := addr.Price()

		assert.NoError(t, err)
		assert.Equal(t, PriceDomestic, price)
	})

	t.Run("Test Init and Price for International Address", func(t *testing.T) {
		addr := &Address{}
		addr.Init(AddressInternational)

		price, err := addr.Price()

		assert.NoError(t, err)
		assert.Equal(t, PriceInternational, price)
	})

	t.Run("Test Price with Incorrect Address", func(t *testing.T) {
		addr := &Address{}
		addr.Init("Unknown")

		price, err := addr.Price()

		assert.Error(t, err)
		assert.Equal(t, 0.0, price)
	})
}
