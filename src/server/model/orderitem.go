package model

import "github.com/shopspring/decimal"

type OrderItem struct {
	ProductID    string          `json:"product_id" bson:"product_id"`
	Name         string          `json:"name" bson:"name"`
	Quantity     uint            `json:"quantity" bson:"quantity"`
	RetailPrice  decimal.Decimal `json:"retail_price" bson:"retail_price"`
	ThumbnailURL string          `json:"thumbnail_url" bson:"thumbnail_url"`
}
