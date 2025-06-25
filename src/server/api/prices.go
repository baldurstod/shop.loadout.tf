package api

import (
	"fmt"
	"strconv"

	"shop.loadout.tf/src/server/databases"
	"shop.loadout.tf/src/server/databases/printfuldb"
	"shop.loadout.tf/src/server/model"
)

func UpdateProductPrice(productId string, currency string) (*model.RetailPrice, error) {
	product, err := databases.GetProduct(productId)
	if err != nil {
		return nil, fmt.Errorf("error while getting product %s: %w", productId, err)
	}

	variantId, err := strconv.Atoi(product.ExternalID1)
	if err != nil {
		return nil, fmt.Errorf("error while converting ExternalID1: %w", err)
	}

	pfVariant, _, err := printfuldb.FindVariant(variantId)
	if err != nil {
		return nil, fmt.Errorf("error while finding variant %d: %w", variantId, err)
	}

	technique, placements, err := product.GetPlacementList()
	if err != nil {
		return nil, fmt.Errorf("error while getting placement list for product variant %s: %w", product.ID, err)
	}

	price, err := computeProductPrice(pfVariant.CatalogProductID, variantId, technique, placements, currency)
	if err != nil {
		return nil, err
	}

	retailPrice, err := databases.SetRetailPrice(product.ID, currency, price)
	if err != nil {
		return nil, err
	}

	return retailPrice, nil
}
