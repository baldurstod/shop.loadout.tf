package model

type TaxInfo struct {
	Required        bool    `json:"required" bson:"required"`
	Rate            float64 `json:"rate" bson:"rate"`
	ShippingTaxable bool    `json:"shipping_taxable" bson:"shipping_taxable"`
}
