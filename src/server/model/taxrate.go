package model

import (
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaxRate struct {
	ID          primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	CountryCode string               `json:"country_code" bson:"country_code"`
	StateCode   string               `json:"state_code" bson:"state_code"`
	PostalCode  string               `json:"postal_code" bson:"postal_code"`
	City        string               `json:"city" bson:"city"`
	Rate        primitive.Decimal128 `json:"rate" bson:"rate"`
}

func NewTaxRate(countryCode string, stateCode string, postalCode string, city string, rate decimal.Decimal) *TaxRate {
	r, _ := primitive.ParseDecimal128(rate.String())

	taxRate := TaxRate{
		CountryCode: countryCode,
		StateCode:   stateCode,
		PostalCode:  postalCode,
		City:        city,
		Rate:        r,
	}

	return &taxRate
}
