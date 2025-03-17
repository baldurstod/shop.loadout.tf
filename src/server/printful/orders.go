package printful

import (
	"errors"

	printfulsdk "github.com/baldurstod/go-printful-sdk"
	"github.com/baldurstod/go-printful-sdk/model"
)

func CalculateShippingRates(recipient model.ShippingRatesAddress, items []model.CatalogOrWarehouseShippingRateItem, currency string, locale string) ([]model.ShippingRate, error) {
	shippingRates, err := printfulClient.CalculateShippingRates(recipient, items, printfulsdk.WithCurrency(currency), printfulsdk.WithLanguage(locale))
	if err != nil {
		return nil, errors.New("unable to get printful response")
	}

	return shippingRates, nil
}
