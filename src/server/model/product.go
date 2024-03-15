package model

type Product struct {
	Id                string        `json:"id" bson:"_id"`
	Name              string        `json:"name" bson:"name"`
	ProductName       string        `json:"product_name" bson:"product_name"`
	ThumbnailUrl      string        `json:"thumbnail_url" bson:"thumbnail_url"`
	Description       string        `json:"description" bson:"description"`
	IsIgnored         bool          `json:"is_ignored" bson:"is_ignored"`
	DateCreated       int64         `json:"date_created" bson:"date_created"`
	DateModified      int64         `json:"date_modified" bson:"date_modified"`
	RetailPrice       float64       `json:"retail_price" bson:"retail_price"`
	Currency          string        `json:"currency" bson:"currency"`
	Files             []File        `json:"files" bson:"files"`
	VariantIds        []string      `json:"variant_ids" bson:"variant_ids"`
	ExternalVariantId int64         `json:"external_variant_id" bson:"external_variant_id"`
	HasMockupPictures bool          `json:"has_mockup_pictures" bson:"has_mockup_pictures"`
	Options           []Option      `json:"options" bson:"options"`
	Variants          []interface{} `json:"variants" bson:"variants"`
	Status            string        `json:"status" bson:"status"`
}
