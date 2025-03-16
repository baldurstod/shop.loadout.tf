package printfulapi

import (
	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"shop.loadout.tf/src/server/mongo/printfuldb"
)

func GetProducts() ([]printfulmodel.Product, error) {
	products, err := printfuldb.FindProducts()

	if err != nil {
		return nil, err
	}

	return products, nil
}
