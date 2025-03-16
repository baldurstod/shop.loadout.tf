package printfulapi

import (
	"errors"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"shop.loadout.tf/src/server/mongo/printfuldb"
)

func GetVariants(productID int) ([]printfulmodel.Variant, error) {
	variants, _, err := printfuldb.FindVariants(productID)
	if err == nil {
		return variants, nil
	}

	return nil, errors.New("unable to find variants")
}
