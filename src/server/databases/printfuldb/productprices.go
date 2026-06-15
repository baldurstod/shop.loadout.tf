package printfuldb

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
)

func FindProductPrices(productID int, currency string) (*printfulmodel.ProductPrices, bool, error) {
	if printfulDb == nil {
		return nil, false, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT product_prices, last_updated FROM products_prices WHERE product_id = $1 AND currency = $2;`
	row := printfulDb.QueryRow(query, productID, currency)

	var productPrices string
	var lastUpdated int64

	err := row.Scan(&productPrices, &lastUpdated)
	if err != nil {
		return nil, false, fmt.Errorf("failed to scan row in FindProductPrices: <%w>", err)
	}

	prices := printfulmodel.ProductPrices{}
	if err = json.Unmarshal([]byte(productPrices), &prices); err != nil {
		return nil, false, err
	}

	return &prices, time.Now().Unix()-lastUpdated > cacheMaxAge, nil
}
