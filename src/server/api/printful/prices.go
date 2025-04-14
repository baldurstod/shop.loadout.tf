package printfulapi

import (
	"errors"
	"fmt"
	"strconv"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"shop.loadout.tf/src/server/databases/printfuldb"
)

func GetProductPrices(productID int, currency string, markup float64) (*printfulmodel.ProductPrices, error) {
	productPrices, _, err := printfuldb.FindProductPrices(productID, currency)
	if err != nil {
		return nil, errors.New("unable to find product prices")
	}

	for i := range productPrices.Product.Placements {
		placement := &productPrices.Product.Placements[i]

		placement.Price, err = applyMarkup(placement.Price, markup) //(1 + printfulConfig.Markup/100)
		if err != nil {
			return nil, errors.New("failed to format product price")
		}

		placement.DiscountedPrice, err = applyMarkup(placement.DiscountedPrice, markup) //(1 + printfulConfig.Markup/100)
		if err != nil {
			return nil, errors.New("failed to format product price")
		}
	}

	for i := range productPrices.Variants {
		variant := &productPrices.Variants[i]
		for j := range variant.Techniques {
			technique := &variant.Techniques[j]
			technique.Price, err = applyMarkup(technique.Price, markup) //(1 + printfulConfig.Markup/100)
			if err != nil {
				return nil, errors.New("failed to format product price")
			}

			technique.DiscountedPrice, err = applyMarkup(technique.DiscountedPrice, markup) //(1 + printfulConfig.Markup/100)
			if err != nil {
				return nil, errors.New("failed to format product price")
			}
		}
	}

	return productPrices, nil
}

func applyMarkup(price string, pct float64) (string, error) {
	p, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return "", err
	}

	p *= (1 + pct*0.01)
	return strconv.FormatFloat(p, 'f', 2, 64), nil
}

type VariantPrice struct {
	ID        int    `json:"id"`
	Technique string `json:"technique"`
	Price     string `json:"price"`
}

func GetVariantsPrices(currency string, markup float64) ([]VariantPrice, error) {
	variantPrices := make([]VariantPrice, 0, 100)

	productPrices, err := printfuldb.FindProductsPrices(currency)
	if err != nil {
		return nil, err
	}

	for _, productPrices := range productPrices {
		for _, variantPrice := range productPrices.Variants {
			for _, technique := range variantPrice.Techniques {

				techniquePrice, err := applyMarkup(technique.Price, markup)
				if err != nil {
					return nil, fmt.Errorf("failed to format technique price %s", technique.Price)
				}

				variantPrices = append(variantPrices, VariantPrice{
					ID:        variantPrice.ID,
					Technique: technique.TechniqueKey,
					Price:     techniquePrice,
				})
			}
		}
	}

	return variantPrices, nil
}
