package model

import (
	"github.com/greatcloak/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Variant struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	ThumbnailURL string             `json:"thumbnail_url" bson:"thumbnail_url"`
	RetailPrice  decimal.Decimal    `json:"retail_price" bson:"retail_price"`
	Currency     string             `json:"currency" bson:"currency"`
	Options      []Option           `json:"options" bson:"options"`
}

func NewVariant(product *Product) Variant {
	return Variant{
		ID:           product.ID,
		Name:         product.Name,
		ThumbnailURL: product.ThumbnailURL,
		//RetailPrice:  product.RetailPrice,
		//Currency: product.Currency,
		Options: product.Options,
	}
}
