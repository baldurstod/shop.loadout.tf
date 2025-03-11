package model

import (
	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"github.com/shopspring/decimal"
)

type Order struct {
	ID                 string                       `json:"id" bson:"id"`
	Currency           string                       `json:"currency" bson:"currency"`
	DateCreated        int64                        `json:"date_created" bson:"date_created"`
	DateUpdated        int64                        `json:"date_updated" bson:"date_updated"`
	ShippingAddress    Address                      `json:"shipping_address" bson:"shipping_address"`
	BillingAddress     Address                      `json:"billing_address" bson:"billing_address"`
	SameBillingAddress bool                         `json:"same_billing_address" bson:"same_billing_address"`
	Items              []OrderItem                  `json:"items" bson:"items"`
	ShippingInfos      []printfulmodel.ShippingRate `json:"shipping_infos" bson:"shipping_infos"`
	TaxInfo            TaxInfo                      `json:"tax_info" bson:"tax_info"`
	ShippingMethod     string                       `json:"shipping_method" bson:"shipping_method"`
	PrintfulOrderID    string                       `json:"printful_order_id" bson:"printful_order_id"`
	PaypalOrderID      string                       `json:"paypal_order_id" bson:"paypal_order_id"`
	Status             string                       `json:"status" bson:"status"`
}

func NewOrder() Order {
	return Order{ShippingInfos: make([]printfulmodel.ShippingRate, 0), SameBillingAddress: true}
}

func (order *Order) GetShippingInfo(shippingMethod string) *printfulmodel.ShippingRate {
	for _, shippingRate := range order.ShippingInfos {
		if shippingMethod == shippingRate.Shipping {
			return &shippingRate
		}
	}
	return nil
}

func (order *Order) GetItemsPrice() *decimal.Decimal {
	price := decimal.Decimal{}
	for _, item := range order.Items {
		price = price.Add(decimal.NewFromInt(int64(item.Quantity)).Mul(item.GetRetailPrice()))
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
