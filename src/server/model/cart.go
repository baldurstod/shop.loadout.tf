package model

import "shop.loadout.tf/src/server/constants"

type Cart struct {
	Currency string `json:"currency" bson:"currency"`
	//Items []CartProduct `json:"products" bson:"products"`
	Items map[string]uint `json:"items" bson:"items"`
}

func NewCart() Cart {
	return Cart{Currency: constants.DEFAULT_CURRENCY, Items: make(map[string]uint)}
}

func (cart Cart) SetQuantity(productID string, quantity uint) {
	if quantity == 0 {
		delete(cart.Items, productID)
	} else {
		cart.Items[productID] = quantity
	}
}

func (cart *Cart) AddQuantity(productID string, quantity uint) {
	cart.Items[productID] = quantity + cart.Items[productID]
}

func (cart *Cart) RemoveProduct(productID string) {
	delete(cart.Items, productID)
}
func (cart *Cart) Clear() {
	cart.Items = make(map[string]uint)
}
func (cart *Cart) TotalQuantity() uint {
	var qty uint
	for _, q := range cart.Items {
		qty += q
	}
	return qty
}
