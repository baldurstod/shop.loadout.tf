package printful

import (
	"errors"
	"fmt"
	"log"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"shop.loadout.tf/src/server/databases/printfuldb"
)

func RefreshAllProducts(currency string, useCache bool) error {
	products, err := printfulClient.GetCatalogProducts()
	if err != nil {
		return errors.New("unable to get catalog products")
	}

	for _, product := range products {
		printfuldb.InsertProduct(&product)
	}

	for _, product := range products {
		if err = refreshVariants(product.ID, product.VariantCount, useCache); err != nil {
			log.Println("Error while refreshing product variants", product.ID, err)
		}

		if err = refreshPrices(product.ID, currency, useCache); err != nil {
			log.Println("Error while refreshing product prices", product.ID, err)
		}

		if err = refreshTemplates(product.ID, useCache); err != nil {
			log.Println("Error while refreshing product templates", product.ID, err)
		}

		if err = refreshStyles(product.ID, useCache); err != nil {
			log.Println("Error while refreshing product styles", product.ID, err)
		}
	}

	return nil
}

func refreshVariants(productID int, count int, useCache bool) error {
	//log.Println("Refreshing variants for product", productID)

	var variants []printfulmodel.Variant
	outdated := true
	var err error

	if useCache {
		variants, outdated, err = printfuldb.FindVariants(productID)
		if err != nil || len(variants) != count {
			outdated = true
		}
	}

	if outdated {
		log.Println("Variants for product", productID, "are outdated, refreshing")
		variants, err = printfulClient.GetCatalogVariants(productID)
		if err != nil {
			//log.Println("Error while getting product variants", productID, err)
			return fmt.Errorf("error while refreshing variants: %w", err)
		} else {

			variantIDs := make([]int, 0, 20)

			for _, variant := range variants {
				variantIDs = append(variantIDs, variant.ID)
				if err = printfuldb.InsertVariant(&variant); err != nil {
					return fmt.Errorf("error while refreshing variants: %w", err)
				}
			}

			if err = printfuldb.UpdateProductVariantIds(productID, variantIDs); err != nil {
				return fmt.Errorf("error while refreshing variants: %w", err)
			}
		}
	}
	return nil
}

func refreshPrices(productID int, currency string, useCache bool) error {
	var prices *printfulmodel.ProductPrices
	outdated := true
	var err error

	if useCache {
		_, outdated, err = printfuldb.FindProductPrices(productID, currency)
		if err != nil {
			outdated = true
		}
	}

	if outdated {
		log.Println("Prices for product", productID, "currency", currency, "are outdated, refreshing")
		prices, err = printfulClient.GetProductPrices(productID)
		if err != nil {
			return fmt.Errorf("error while refreshing prices: %w", err)
		} else {
			printfuldb.InsertProductPrices(prices)
		}
	}

	return nil
}

func refreshTemplates(productID int, useCache bool) error {
	var templates []printfulmodel.MockupTemplates
	outdated := true
	var err error

	if useCache {
		_, outdated, err = printfuldb.FindMockupTemplates(productID)
		if err != nil {
			outdated = true
		}
	}

	if outdated {
		log.Println("Templates for product", productID, "are outdated, refreshing")
		templates, err = printfulClient.GetMockupTemplates(productID)
		if err != nil {
			return fmt.Errorf("error while refreshing templates: %w", err)
		} else {
			printfuldb.InsertMockupTemplates(productID, templates)
		}
	}

	return nil
}

func refreshStyles(productID int, useCache bool) error {
	var styles []printfulmodel.MockupStyles
	outdated := true
	var err error

	if useCache {
		_, outdated, err = printfuldb.FindMockupStyles(productID)
		if err != nil {
			outdated = true
		}
	}

	if outdated {
		log.Println("Styles for product", productID, "are outdated, refreshing")
		styles, err = printfulClient.GetMockupStyles(productID)
		if err != nil {
			return fmt.Errorf("error while refreshing styles: %w", err)
		} else {
			printfuldb.InsertMockupStyles(productID, styles)
		}
	}

	return nil
}
