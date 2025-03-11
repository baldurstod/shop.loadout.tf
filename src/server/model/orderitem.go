package model

import (
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ProductID    string               `json:"product_id" bson:"product_id"`
	Name         string               `json:"name" bson:"name"`
	Quantity     uint                 `json:"quantity" bson:"quantity"`
	RetailPrice  primitive.Decimal128 `json:"retail_price" bson:"retail_price"`
	ThumbnailURL string               `json:"thumbnail_url" bson:"thumbnail_url"`
}

func (p *OrderItem) GetRetailPrice() decimal.Decimal {
	price, _ := decimal.NewFromString(p.RetailPrice.String())
	return price
}
