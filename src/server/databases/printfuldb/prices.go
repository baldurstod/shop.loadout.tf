package printfuldb

import (
	"encoding/json"
	"errors"
	"fmt"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
)

func FindProductsPrices(currency string) ([]printfulmodel.ProductPrices, error) {
	if printfulDb == nil {
		return nil, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT product_id, product_prices, last_updated FROM products_prices WHERE currency = $1;`
	res, err := printfulDb.Query(query, currency)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query "+query+"in FindProductsPrices: <%w>", err)
	}
	defer res.Close()

	prices := make([]printfulmodel.ProductPrices, 0, 400)
	for res.Next() {
		var productID int
		var productPrices string
		var lastUpdated int64

		err = res.Scan(&productID, &productPrices, &lastUpdated)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row in FindCategories: <%w>", err)
		}

		price := printfulmodel.ProductPrices{}
		if err = json.Unmarshal([]byte(productPrices), &price); err != nil {
			return nil, err
		}

		//category := printfulmodel.ProductPrices{ID: id, ParentID: parent_id, ImageURL: image_url, Title: title}

		prices = append(prices, price)
	}

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("failed to get next row in FindCategories: <%w>", err)
	}

	return prices, nil
}

/*
func FindProductsPrices(currency string) ([]*printfulmodel.ProductPrices, error) {
	ctx, cancel := context.WithTimeout(context.Background(), databases.MongoTimeout)
	defer cancel()

	filter := bson.D{{Key: "currency", Value: currency}}

	cursor, err := pfProductsPricesCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	prices := make([]*printfulmodel.ProductPrices, 0, 400)
	for cursor.Next(context.TODO()) {
		doc := MongoProductPrices{}
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}

		prices = append(prices, &doc.ProductPrices)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return prices, nil
}
*/
