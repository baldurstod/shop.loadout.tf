package model

import (
	"github.com/greatcloak/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RetailPrice struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	ProductID   primitive.ObjectID `json:"product_id" bson:"product_id"`
	Currency    string             `json:"currency" bson:"currency"`
	RetailPrice decimal.Decimal    `json:"retail_price" bson:"retail_price"`
}
