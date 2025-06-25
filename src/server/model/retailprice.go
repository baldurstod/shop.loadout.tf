package model

import (
	"time"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RetailPrice struct {
	ID          primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	ProductID   string               `json:"product_id" bson:"product_id"`
	Currency    string               `json:"currency" bson:"currency"`
	RetailPrice primitive.Decimal128 `json:"retail_price" bson:"retail_price"`
	DateUpdated int64                `json:"date_updated" bson:"date_updated"`
}

func NewRetailPrice(productId string, currency string, price decimal.Decimal) *RetailPrice {
	p := RetailPrice{
		ProductID:   productId,
		Currency:    currency,
		DateUpdated: time.Now().Unix(),
	}

	p.SetPrice((price))

	return &p
}

func (p *RetailPrice) SetPrice(retailPrice decimal.Decimal) {
	p.RetailPrice, _ = primitive.ParseDecimal128(retailPrice.String())
}

func (p *RetailPrice) GetPrice() decimal.Decimal {
	price, _ := decimal.NewFromString(p.RetailPrice.String())
	return price
}
