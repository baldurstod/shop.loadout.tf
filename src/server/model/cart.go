package model

import "shop.loadout.tf/src/server/constants"

type Cart struct {
	Currency string `json:"currency" bson:"currency"`
	//Items []CartProduct `json:"products" bson:"products"`
	Items map[string]int64 `json:"items" bson:"items"`
}

func NewCart() Cart {
	return Cart{Currency: constants.DEFAULT_CURRENCY, Items: make(map[string]int64)}
}

func (cart Cart) SetQuantity(productID string, quantity int64) {
	if quantity == 0 {
		delete(cart.Items, productID)
	} else {
		cart.Items[productID] = quantity
	}
}

func (cart *Cart) AddQuantity(productID string, quantity int64) {
	if cart.Items == nil {
		cart.Items = make(map[string]int64)
	}

	cart.Items[productID] = quantity + cart.Items[productID]
}

func (cart *Cart) RemoveProduct(productID string) {
	delete(cart.Items, productID)
}
func (cart *Cart) Clear() {
	cart.Items = make(map[string]int64)
}
func (cart *Cart) TotalQuantity() int64 {
	if cart.Items == nil {
		return 0
	}

	var qty int64
	for _, q := range cart.Items {
		qty += q
	}
	return qty
}
