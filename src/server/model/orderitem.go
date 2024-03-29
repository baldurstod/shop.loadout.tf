package model

type OrderItem struct {
	ProductID    string  `json:"product_id" bson:"product_id"`
	Name         string  `json:"name" bson:"name"`
	Quantity     uint    `json:"quantity" bson:"quantity"`
	RetailPrice  float64 `json:"retail_price" bson:"retail_price"`
	ThumbnailUrl string  `json:"thumbnail_url" bson:"thumbnail_url"`
}
