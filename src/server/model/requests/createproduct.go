package requests

import (
	"errors"
)

type CreateProductRequest struct {
	VariantID  int                             `json:"variant_id"`
	Name       string                          `json:"name"`
	Placements []CreateProductRequestPlacement `json:"placements"`
}

type CreateProductRequestPlacement struct {
	Name      string `json:"name"`
	Technique string `json:"technique"`
}

func (request CreateProductRequest) CheckParams() error {
	if request.VariantID == 0 {
		return errors.New("Invalid variant id")
	}

	if request.Name == "" {
		return errors.New("Invalid name")
	}
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
