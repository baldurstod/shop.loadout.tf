package model

import (
	"fmt"
	"slices"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"github.com/mitchellh/mapstructure"
)

type Product struct {
	ID           string         `json:"id" bson:"id"`
	Name         string         `json:"name" bson:"name"`
	ProductName  string         `json:"product_name" bson:"product_name"`
	ThumbnailURL string         `json:"thumbnail_url" bson:"thumbnail_url"`
	Description  string         `json:"description" bson:"description"`
	IsIgnored    bool           `json:"is_ignored" bson:"is_ignored"`
	DateCreated  int64          `json:"date_created" bson:"date_created"`
	DateUpdated  int64          `json:"date_updated" bson:"date_updated"`
	Files        []File         `json:"files" bson:"files"`
	VariantIDs   []string       `json:"variant_ids" bson:"variant_ids"`
	ExternalID1  string         `json:"external_id_1" bson:"external_id_1"`
	ExternalID2  string         `json:"external_id_2" bson:"external_id_2"`
	ExtraData    map[string]any `json:"extra_data" bson:"extra_data"`
	Options      []Option       `json:"options" bson:"options"`
	Variants     []Variant      `json:"variants" bson:"variants"`
	Status       string         `json:"status" bson:"status"`
}

func NewProduct() Product {
	return Product{
		Files:      []File{},
		VariantIDs: []string{},
		Options:    []Option{},
		Variants:   []Variant{},
	}
}

func (product *Product) AddOption(name string, optionType string, optionValue string) {
	product.Options = append(product.Options, Option{
		Name:  name,
		Type:  optionType,
		Value: optionValue,
	})
}

func (product *Product) SetFile(fileType string, url string, thumbURL string) {
	product.Files = slices.DeleteFunc(product.Files, func(n File) bool {
		return n.Type == fileType
	})

	product.Files = append(product.Files, File{
		Type:         fileType,
		URL:          url,
		ThumbnailURL: thumbURL,
	})
}

func (product *Product) AddVariant(variant Variant) {
	product.Variants = append(product.Variants, variant)
}

func (product *Product) GetPlacementList() (string, printfulmodel.PlacementsList, error) {
	productExtraData := ProductExtraData{}
	err := mapstructure.Decode(product.ExtraData, &productExtraData)
	if err != nil {
		return "", nil, fmt.Errorf("error while decoding product extra data for product %s: %w", product.ID, err)
	}

	placementsList := make(printfulmodel.PlacementsList, len(productExtraData.Printful.Placements))

	for i, placement := range productExtraData.Printful.Placements {
		placementsList[i] = printfulmodel.Placement{
			Placement:     placement.Placement,
			Technique:     placement.Technique,
			PrintAreaType: "simple", //TODO: variable ?
			Layers: []printfulmodel.Layer{{
				Type: "file", //TODO: variable ?
				Url:  placement.ImageURL,
				//LayerOptions
				//LayerPosition
			}},
		}
	}

	return productExtraData.Printful.Technique, placementsList, nil
}
