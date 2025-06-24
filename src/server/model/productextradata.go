package model

type ProductExtraData struct {
	Printful ProductExtraDataPrintful `json:"printful" bson:"printful" mapstructure:"printful"`
}

type ProductExtraDataPrintful struct {
	Technique  string                      `json:"technique" bson:"technique" mapstructure:"technique"`
	Placements []ProductExtraDataPlacement `json:"placements" bson:"placements" mapstructure:"placements"`
}

type ProductExtraDataPlacement struct {
	Placement   string `json:"placement" bson:"placement" mapstructure:"placement"`
	Technique   string `json:"technique" bson:"technique" mapstructure:"technique"`
	Orientation string `json:"orientation" bson:"orientation" mapstructure:"orientation"`
	ImageURL    string `json:"image_url" bson:"image_url" mapstructure:"image_url"`
	ThumbURL    string `json:"thumb_url" bson:"thumb_url" mapstructure:"thumb_url"`
}
