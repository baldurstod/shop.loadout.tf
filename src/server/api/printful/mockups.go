package printfulapi

import (
	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"shop.loadout.tf/src/server/databases/printfuldb"
)

func GetMockupTemplates(productID int) ([]printfulmodel.MockupTemplates, error) {
	templates, _, err := printfuldb.FindMockupTemplates(productID)

	if err != nil {
		return nil, err
	}

	return templates, nil
}

func GetMockupStyles(productID int) ([]printfulmodel.MockupStyles, error) {
	styles, _, err := printfuldb.FindMockupStyles(productID)

	if err != nil {
		return nil, err
	}

	return styles, nil
}
