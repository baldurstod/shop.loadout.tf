package model

import (
	"time"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID                 string                       `json:"id" bson:"id"`
	Currency           string                       `json:"currency" bson:"currency"`
	ShippingAddress    Address                      `json:"shipping_address" bson:"shipping_address"`
	BillingAddress     Address                      `json:"billing_address" bson:"billing_address"`
	SameBillingAddress bool                         `json:"same_billing_address" bson:"same_billing_address"`
	Items              []OrderItem                  `json:"items" bson:"items"`
	ShippingInfos      []printfulmodel.ShippingRate `json:"shipping_infos" bson:"shipping_infos"`
	TaxInfo            TaxInfo                      `json:"tax_info" bson:"tax_info"`
	PercentDiscount    primitive.Decimal128         `json:"percent_discount" bson:"percent_discount"`
	PriceDiscount      primitive.Decimal128         `json:"price_discount" bson:"price_discount"`
	ShippingMethod     string                       `json:"shipping_method" bson:"shipping_method"`
	PrintfulOrderID    string                       `json:"printful_order_id" bson:"printful_order_id"`
	PaypalOrderID      string                       `json:"paypal_order_id" bson:"paypal_order_id"`
	Status             string                       `json:"status" bson:"status"`
	DateCreated        time.Time                    `json:"date_created" bson:"date_created"`
	DateUpdated        time.Time                    `json:"date_updated" bson:"date_updated"`
}

func NewOrder() Order {
	percent, _ := primitive.ParseDecimal128("0.1")
	return Order{ShippingInfos: make([]printfulmodel.ShippingRate, 0), SameBillingAddress: true, PercentDiscount: percent, Items: make([]OrderItem, 0)}
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
	percentOff, _ := decimal.NewFromString(order.PercentDiscount.String())
	priceDiscount, _ := decimal.NewFromString(order.PriceDiscount.String())
	price := order.GetItemsPrice().Mul(decimal.NewFromInt(1).Sub(percentOff)).Sub(priceDiscount).Add(*order.GetShippingPrice()).Add(*order.GetTaxPrice())

	price = price.Round(2)
	return &price
}
