package printfulapi

import (
	"errors"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"shop.loadout.tf/src/server/databases/printfuldb"
)

func GetProducts() ([]printfulmodel.Product, error) {
	products, err := printfuldb.FindProducts()

	if err != nil {
		return nil, err
	}

	return products, nil
}

func GetProduct(productID int) (*printfulmodel.Product, error) {
	product, _, err := printfuldb.FindProduct(productID)
	if err == nil {
		return product, nil
	}

	return nil, errors.New("unable to find product")
}
