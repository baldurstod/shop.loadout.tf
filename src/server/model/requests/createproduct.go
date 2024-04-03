package requests

import (
	"encoding/base64"
	"errors"
	"image"
	"image/png"
	"strings"
)

type CreateProductRequest struct {
	VariantID    int    `mapstructure:"variant_id"`
	Name         string `mapstructure:"name"`
	Type         string `mapstructure:"type"`
	Image        string `mapstructure:"image"`
	DecodedImage image.Image
}

func (request CreateProductRequest) CheckParams() error {
	if request.VariantID == 0 {
		return errors.New("Invalid variant id")
	}

	if request.Name == "" {
		return errors.New("Invalid name")
	}

	if request.Type == "" {
		return errors.New("Invalid type")
	}

	if request.Image == "" {
		return errors.New("Invalid image")
	}

	b64data := request.Image[strings.IndexByte(request.Image, ',')+1:] // Remove data:image/png;base64,
	img, err := png.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64data)))
	if err != nil {
		return errors.New("Invalid image")
	}

	request.DecodedImage = img

	return nil
}
