package model

import (
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RetailPrice struct {
	ID          primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	ProductID   primitive.ObjectID   `json:"product_id" bson:"product_id"`
	Currency    string               `json:"currency" bson:"currency"`
	RetailPrice primitive.Decimal128 `json:"retail_price" bson:"retail_price"`
	DateUpdated int64                `json:"date_updated" bson:"date_updated"`
}

func (p *RetailPrice) SetRetailPrice(retailPrice decimal.Decimal) {
	p.RetailPrice, _ = primitive.ParseDecimal128(retailPrice.String())
}

func (p *RetailPrice) GetRetailPrice() decimal.Decimal {
	price, _ := decimal.NewFromString(p.RetailPrice.String())
	return price
}
