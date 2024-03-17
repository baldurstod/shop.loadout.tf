package model

type Cart struct {
	Currency string        `json:"currency" bson:"currency"`
	//Items []CartProduct `json:"products" bson:"products"`
	Items map[string]uint `json:"items" bson:"items"`
}

func NewCart() Cart {
	return Cart{Items: make(map[string]uint)}
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
