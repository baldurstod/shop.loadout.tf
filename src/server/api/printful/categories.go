package printfulapi

import (
	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"shop.loadout.tf/src/server/databases/printfuldb"
)

func GetCategories(language string) ([]printfulmodel.Category, error) {
	categories, err := printfuldb.GetCategories(language)

	if err != nil {
		return nil, err
	}

	return categories, nil
}
