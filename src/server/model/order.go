package model

import (
	"fmt"
	"time"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"github.com/shopspring/decimal"
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
	PercentDiscount    decimal.Decimal              `json:"percent_discount" bson:"percent_discount"`
	PriceDiscount      decimal.Decimal              `json:"price_discount" bson:"price_discount"`
	ShippingMethod     string                       `json:"shipping_method" bson:"shipping_method"`
	ItemsPrice         decimal.Decimal              `json:"items_price"`
	DiscountPrice      decimal.Decimal              `json:"discount_price"`
	ShippingPrice      decimal.Decimal              `json:"shipping_price"`
	TaxPrice           decimal.Decimal              `json:"tax_price"`
	TotalPrice         decimal.Decimal              `json:"total_price"`
	PrintfulOrderID    string                       `json:"printful_order_id" bson:"printful_order_id"`
	PaypalOrderID      string                       `json:"paypal_order_id" bson:"paypal_order_id"`
	Status             string                       `json:"status" bson:"status"`
	DateCreated        time.Time                    `json:"date_created" bson:"date_created"`
	DateUpdated        time.Time                    `json:"date_updated" bson:"date_updated"`
}

func NewOrder() Order {
	percent := decimal.NewFromFloat32(0.1)
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

func (order *Order) GetShippingPrice() (*decimal.Decimal, error) {
	if order.ShippingMethod == "" {
		return &decimal.Decimal{}, nil
	}
	shippingInfo := order.GetShippingInfo(order.ShippingMethod)
	if shippingInfo == nil {
		return nil, fmt.Errorf("shipping info for method %s not found", order.ShippingMethod)
	}
	price, err := decimal.NewFromString(shippingInfo.Rate)
	if err != nil {
		return nil, fmt.Errorf("unable to convert shipping rate %s to decimal: %w", shippingInfo.Rate, err)
	}
	price = price.Round(2)
	return &price, nil
}

func (order *Order) GetTaxPrice() (*decimal.Decimal, error) {
	taxRate := decimal.NewFromFloat(order.TaxInfo.Rate)
	price := order.GetItemsPrice().Mul(taxRate)

	if order.TaxInfo.ShippingTaxable {
		shippingPrice, err := order.GetShippingPrice()
		if err != nil {
			return nil, err
		}
		price = price.Add(shippingPrice.Mul(taxRate))
	}

	price = price.Round(2)
	return &price, nil
}

func (order *Order) GetTotalPrice() (*decimal.Decimal, error) {
	percentOff, _ := decimal.NewFromString(order.PercentDiscount.String())
	priceDiscount, _ := decimal.NewFromString(order.PriceDiscount.String())

	shippingPrice, err := order.GetShippingPrice()
	if err != nil {
		return nil, err
	}

	taxPrice, err := order.GetTaxPrice()
	if err != nil {
		return nil, err
	}

	price := order.GetItemsPrice().Mul(decimal.NewFromInt(1).Sub(percentOff)).Sub(priceDiscount).Add(*shippingPrice).Add(*taxPrice)

	price = price.Round(2)
	return &price, nil
}
