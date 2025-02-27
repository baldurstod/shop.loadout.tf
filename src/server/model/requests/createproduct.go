package requests

import "image"

type CreateProductRequest struct {
	ProductID int `json:"product_id"  mapstructure:"product_id"`
	VariantID int `json:"variant_id"  mapstructure:"variant_id"`
	//Name       string                          `json:"name"  mapstructure:"name"`
	Placements []CreateProductRequestPlacement `json:"placements"  mapstructure:"placements"`
}

type CreateProductRequestPlacement struct {
	Placement    string `json:"placement"  mapstructure:"placement"`
	Technique    string `json:"technique"  mapstructure:"technique"`
	Image        string `json:"image"  mapstructure:"image"`
	Orientation  string `json:"orientation"  mapstructure:"orientation"`
	DecodedImage image.Image
}
