package requests

import (
	"errors"
)

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

func (request CreateProductRequest) CheckParams() error {
	if request.VariantID == 0 {
		return errors.New("invalid variant id")
	}
	/*
		if request.Name == "" {
			return errors.New("invalid name")
		}
	*/
	/*
		if request.Image == "" {
			return errors.New("Invalid image")
		}

		b64data := request.Image[strings.IndexByte(request.Image, ',')+1:] // Remove data:image/png;base64,
		img, err := png.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64data)))
		if err != nil {
			return errors.New("Invalid image")
		}

		request.DecodedImage = img
	*/

	return nil
}
