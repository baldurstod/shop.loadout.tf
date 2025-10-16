package printful

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/baldurstod/go-printful-api-model/responses"
	"github.com/baldurstod/go-printful-api-model/schemas"
	printfulsdk "github.com/baldurstod/go-printful-sdk"
	"github.com/baldurstod/go-printful-sdk/model"
	"github.com/mitchellh/mapstructure"
)

func CalculateShippingRates(recipient model.ShippingRatesAddress, items []model.CatalogOrWarehouseShippingRateItem, currency string, locale string) ([]model.ShippingRate, error) {
	shippingRates, err := printfulClient.CalculateShippingRates(recipient, items, printfulsdk.WithCurrency(currency), printfulsdk.WithLanguage(locale))
	if err != nil {
		return nil, errors.New("unable to calculate shipping rates")
	}

	return shippingRates, nil
}

func CalculateTaxRate(recipient schemas.TaxAddressInfo) (*schemas.TaxInfo, error) {
	//TODO: this use api v1: find a better solution
	type CalculateTaxRate struct {
		Recipient schemas.TaxAddressInfo `json:"recipient" bson:"recipient" mapstructure:"recipient"`
	}

	ctr := CalculateTaxRate{
		Recipient: recipient,
	}

	body := map[string]any{}
	err := mapstructure.Decode(ctr, &body)
	if err != nil {
		return nil, fmt.Errorf("error while decoding calculate tax rate param %v: %w", ctr, err)
	}

	headers := map[string]string{
		"Authorization": "Bearer " + printfulConfig.AccessToken,
	}

	resp, err := fetchRateLimited("POST", PRINTFUL_TAX_API, "/rates", headers, body)
	if err != nil {
		return nil, fmt.Errorf("unable to calculate tax rate: %w", err)
	}
	defer resp.Body.Close()

	//response := map[string]any{}
	response := responses.TaxRates{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error while decoding calculate tax rate response: %w", err)
	}

	//p := &(response.Result)

	return &response.Result, nil
}

func CreateOrder(externalID string, shipping string, recipient model.Address, orderItems []model.CatalogItem, customization *model.Customization, retailCosts *model.RetailCosts2) (*model.Order, error) {
	opts := make([]printfulsdk.RequestOption, 0, 5)

	if externalID != "" {
		opts = append(opts, printfulsdk.SetOrderExternalID(externalID))
	}

	if shipping != "" {
		opts = append(opts, printfulsdk.SetOrderShippingMethod(shipping))
	}

	if customization != nil {
		opts = append(opts, printfulsdk.SetOrderCustomization(customization))
	}

	if retailCosts != nil {
		opts = append(opts, printfulsdk.SetOrderRetailCosts(retailCosts))
	}

	order, err := printfulClient.CreateOrder(recipient, orderItems, opts...)
	if err != nil {
		return nil, errors.New("unable to create third party order")
	}

	return order, nil
}
