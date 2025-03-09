package model

import (
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RetailPrice struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ProductID   primitive.ObjectID `json:"product_id" bson:"product_id"`
	Currency    string             `json:"currency" bson:"currency"`
	RetailPrice decimal.Decimal    `json:"retail_price" bson:"retail_price"`
	DateUpdated int64              `json:"date_updated" bson:"date_updated"`
}
