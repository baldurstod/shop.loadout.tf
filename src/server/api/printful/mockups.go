package printfulapi

import (
	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"shop.loadout.tf/src/server/mongo/printfuldb"
)

func GetMockupStyles(productID int) ([]printfulmodel.MockupStyles, error) {
	styles, _, err := printfuldb.FindMockupStyles(productID)

	if err != nil {
		return nil, err
	}

	return styles, nil
}
