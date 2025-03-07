package model

import (
	"github.com/baldurstod/go-printful-api-model/schemas"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID                 primitive.ObjectID     `json:"id" bson:"_id"`
	Currency           string                 `json:"currency" bson:"currency"`
	DateCreated        int64                  `json:"date_created" bson:"date_created"`
	DateUpdated        int64                  `json:"date_updated" bson:"date_updated"`
	ShippingAddress    Address                `json:"shipping_address" bson:"shipping_address"`
	BillingAddress     Address                `json:"billing_address" bson:"billing_address"`
	SameBillingAddress bool                   `json:"same_billing_address" bson:"same_billing_address"`
	Items              []OrderItem            `json:"items" bson:"items"`
	ShippingInfos      []schemas.ShippingInfo `json:"shipping_infos" bson:"shipping_infos"`
	TaxInfo            TaxInfo                `json:"tax_info" bson:"tax_info"`
	ShippingMethod     string                 `json:"shipping_method" bson:"shipping_method"`
	PrintfulOrderID    string                 `json:"printful_order_id" bson:"printful_order_id"`
	PaypalOrderID      string                 `json:"paypal_order_id" bson:"paypal_order_id"`
	Status             string                 `json:"status" bson:"status"`
}

func NewOrder() Order {
	return Order{ShippingInfos: make([]schemas.ShippingInfo, 0), SameBillingAddress: true}
}

func (order *Order) GetShippingInfo(shippingMethod string) *schemas.ShippingInfo {
	for _, shippingInfo := range order.ShippingInfos {
		if shippingMethod == shippingInfo.ID {
			return &shippingInfo
		}
	}
	return nil
}

func (order *Order) GetItemsPrice() *decimal.Decimal {
	price := decimal.Decimal{}
	for _, item := range order.Items {
		price = price.Add(decimal.NewFromInt(int64(item.Quantity)).Mul(item.RetailPrice))
	}

	price = price.Round(2)
	return &price
}

func (order *Order) GetShippingPrice() *decimal.Decimal {
	shippingInfo := order.GetShippingInfo(order.ShippingMethod)
	if shippingInfo != nil {
		price, err := decimal.NewFromString(shippingInfo.Rate)
		if err == nil {
			price = price.Round(2)
			return &price
		}
	}

	return &decimal.Decimal{}
}

func (order *Order) GetTaxPrice() *decimal.Decimal {
	taxRate := decimal.NewFromFloat(order.TaxInfo.Rate)
	price := order.GetItemsPrice().Mul(taxRate)

	if order.TaxInfo.ShippingTaxable {
		price = price.Add(order.GetShippingPrice().Mul(taxRate))
	}

	price = price.Round(2)
	return &price
}

func (order *Order) GetTotalPrice() *decimal.Decimal {
	price := order.GetItemsPrice().Add(*order.GetShippingPrice()).Add(*order.GetTaxPrice())

	price = price.Round(2)
	return &price
}
