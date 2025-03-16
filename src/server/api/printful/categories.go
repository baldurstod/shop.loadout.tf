package printfulapi

import (
	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"shop.loadout.tf/src/server/mongo/printfuldb"
)

func GetCategories() ([]printfulmodel.Category, error) {
	categories, err := printfuldb.FindCategories()

	if err != nil {
		return nil, err
	}

	return categories, nil
}
