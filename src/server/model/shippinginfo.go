package model

type ShippingInfo struct {
	ID              string  `json:"id" bson:"id"`
	Name            string  `json:"name" bson:"name"`
	Rate            float64 `json:"rate" bson:"rate"`
	Currency        string  `json:"currency" bson:"currency"`
	MinDeliveryDays uint    `json:"min_delivery_days" bson:"min_delivery_days"`
	MaxDeliveryDays uint    `json:"max_delivery_days" bson:"max_delivery_days"`
}
