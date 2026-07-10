package model

import (
	"github.com/shopspring/decimal"
)

type TaxRate struct {
	CountryCode string          `json:"country_code" bson:"country_code"`
	StateCode   string          `json:"state_code" bson:"state_code"`
	PostalCode  string          `json:"postal_code" bson:"postal_code"`
	City        string          `json:"city" bson:"city"`
	Rate        decimal.Decimal `json:"rate" bson:"rate"`
}

func NewTaxRate(countryCode string, stateCode string, postalCode string, city string, rate decimal.Decimal) *TaxRate {
	//r, _ := primitive.ParseDecimal128(rate.String())

	taxRate := TaxRate{
		CountryCode: countryCode,
		StateCode:   stateCode,
		PostalCode:  postalCode,
		City:        city,
		Rate:        rate,
	}

	return &taxRate
}
