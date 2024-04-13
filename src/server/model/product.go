package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	Name              string             `json:"name" bson:"name"`
	ProductName       string             `json:"product_name" bson:"product_name"`
	ThumbnailURL      string             `json:"thumbnail_url" bson:"thumbnail_url"`
	Description       string             `json:"description" bson:"description"`
	IsIgnored         bool               `json:"is_ignored" bson:"is_ignored"`
	DateCreated       int64              `json:"date_created" bson:"date_created"`
	DateUpdated       int64              `json:"date_updated" bson:"date_updated"`
	RetailPrice       float64            `json:"retail_price" bson:"retail_price"`
	Currency          string             `json:"currency" bson:"currency"`
	Files             []File             `json:"files" bson:"files"`
	VariantIds        []string           `json:"variant_ids" bson:"variant_ids"`
	ExternalVariantId int64              `json:"external_variant_id" bson:"external_variant_id"`
	HasMockupPictures bool               `json:"has_mockup_pictures" bson:"has_mockup_pictures"`
	Options           []Option           `json:"options" bson:"options"`
	Variants          []interface{}      `json:"variants" bson:"variants"`
	Status            string             `json:"status" bson:"status"`
}

func NewProduct() Product {
	return Product{}
}

func (product Product) AddOption(name string, optionType string, optionValue string) {
	product.Options = append(product.Options, Option{
		Name: name,
		Type: optionType,
		Value: optionValue,
	});
}

func (product Product) AddFile(fileType string, url string) {
	product.Files = append(product.Files, File{
		Type: fileType,
		URL: url,
	});
}
