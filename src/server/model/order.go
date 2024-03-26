package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID                 primitive.ObjectID      `json:"id" bson:"_id"`
	Currency           string                  `json:"currency" bson:"currency"`
	DateCreated        int64                   `json:"date_created" bson:"date_created"`
	DateUpdated        int64                   `json:"date_updated" bson:"date_updated"`
	ShippingAddress    Address                 `json:"shipping_address" bson:"shipping_address"`
	BillingAddress     Address                 `json:"billing_address" bson:"billing_address"`
	SameBillingAddress bool                    `json:"same_billing_address" bson:"same_billing_address"`
	Items              []OrderItem             `json:"items" bson:"items"`
	ShippingInfos      map[string]ShippingInfo `json:"shipping_infos" bson:"shipping_infos"`
	TaxInfo            TaxInfo                 `json:"tax_info" bson:"tax_info"`
	ShippingMethod     string                  `json:"shipping_method" bson:"shipping_method"`
	PrintfulOrderID    string                  `json:"printful_order_id" bson:"printful_order_id"`
	PaypalOrderID      string                  `json:"paypal_order_id" bson:"paypal_order_id"`
	Status             string                  `json:"status" bson:"status"`
}

func NewOrder() Order {
	return Order{ShippingInfos: make(map[string]ShippingInfo), SameBillingAddress: true}
}
