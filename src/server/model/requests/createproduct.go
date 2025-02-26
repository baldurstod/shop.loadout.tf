package requests

type CreateProductRequest struct {
	VariantID int `json:"variant_id"  mapstructure:"variant_id"`
	//Name       string                          `json:"name"  mapstructure:"name"`
	Placements []CreateProductRequestPlacement `json:"placements"  mapstructure:"placements"`
}

type CreateProductRequestPlacement struct {
	Id          string `json:"id"  mapstructure:"id"`
	Technique   string `json:"technique"  mapstructure:"technique"`
	Image       string `json:"image"  mapstructure:"image"`
	Orientation string `json:"orientation"  mapstructure:"orientation"`
}
