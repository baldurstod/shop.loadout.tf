package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserInfo struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	FirstName       string             `json:"first_name" bson:"first_name"`
	LastName        string             `json:"last_name" bson:"last_name"`
	DateCreated     int64              `json:"date_created" bson:"date_created"`
	DateUpdated     int64              `json:"date_updated" bson:"date_updated"`
	ShippingAddress Address            `json:"shipping_address" bson:"shipping_address"`
	BillingAddress  Address            `json:"billing_address" bson:"billing_address"`
}
