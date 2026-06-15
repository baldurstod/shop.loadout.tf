package printfuldb

import (
	"errors"
	"fmt"
	"log"
	"time"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
)

func FindVariants(productID int) (variants []printfulmodel.Variant, outdated bool, err error) {
	if printfulDb == nil {
		return nil, false, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	outdated = false

	query := `SELECT id, name, catalog_product_id, color, color_code, color_code2, image, size, availability, last_updated FROM variants WHERE catalog_product_id = $1;`
	res, err := printfulDb.Query(query, productID)
	if err != nil {
		return nil, false, fmt.Errorf("failed to execute query "+query+"in FindVariants: <%w>", err)
	}
	defer res.Close()

	variants = make([]printfulmodel.Variant, 0, 20)
	for res.Next() {
		var id int
		var name string
		var catalogProductID int
		var color string
		var colorCode string
		var colorCode2 string
		var image string
		var size string
		var availability string
		var lastUpdated int64

		err = res.Scan(&id, &name, &catalogProductID, &color, &colorCode, &colorCode2, &image, &size, &availability, &lastUpdated)
		if err != nil {
			return nil, false, fmt.Errorf("failed to scan row in FindVariants: <%w>", err)
		}
		variant := printfulmodel.Variant{
			ID:               id,
			Name:             name,
			CatalogProductID: catalogProductID,
			Color:            color,
			ColorCode:        colorCode,
			ColorCode2:       colorCode2,
			Image:            image,
			Size:             size,
			//Availability:     availability,
		}

		if time.Now().Unix()-(lastUpdated) > cacheMaxAge {
			outdated = true
		}

		log.Println("todo: availability")

		variants = append(variants, variant)
	}

	if err := res.Err(); err != nil {
		return nil, false, fmt.Errorf("failed to get next row in FindVariants: <%w>", err)
	}

	return variants, outdated, nil
}

func FindVariant(variantID int) (*printfulmodel.Variant, bool, error) {
	if printfulDb == nil {
		return nil, false, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT id, name, catalog_product_id, color, color_code, color_code2, image, size, availability, last_updated FROM variants WHERE id = $1;`
	row := printfulDb.QueryRow(query, variantID)

	var id int
	var name string
	var catalogProductID int
	var color string
	var colorCode string
	var colorCode2 string
	var image string
	var size string
	var availability string
	var lastUpdated int64

	err := row.Scan(&id, &name, &catalogProductID, &color, &colorCode, &colorCode2, &image, &size, &availability, &lastUpdated)
	if err != nil {
		return nil, false, fmt.Errorf("failed to scan row in FindProduct: <%w>", err)
	}

	variant := printfulmodel.Variant{
		ID:               id,
		Name:             name,
		CatalogProductID: catalogProductID,
		Color:            color,
		ColorCode:        colorCode,
		ColorCode2:       colorCode2,
		Image:            image,
		Size:             size,
		//Availability:     availability,
	}

	return &variant, time.Now().Unix()-lastUpdated > cacheMaxAge, nil
}
