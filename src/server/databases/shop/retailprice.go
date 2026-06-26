package shop

import (
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"shop.loadout.tf/src/server/model"
)

/*
func SetRetailPrice(productID string, currency string, price decimal.Decimal) (*model.RetailPrice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	retailPrice := model.NewRetailPrice(productID, currency, price)

	retailPrice.SetPrice(price)

	/*

		ProductID   primitive.ObjectID `json:"product_id" bson:"product_id"`
		Currency    string             `json:"currency" bson:"currency"`
		RetailPrice decimal.Decimal    `json:"retail_price" bson:"retail_price"`
	* /

	opts := options.Replace().SetUpsert(true)

	filter := bson.D{primitive.E{Key: "product_id", Value: productID}, primitive.E{Key: "currency", Value: currency}}
	_, err := retailPriceCollection.ReplaceOne(ctx, filter, retailPrice, opts)
	if err != nil {
		return nil, err
	}

	return retailPrice, nil
}
*/

func InsertRetailPrice(retailPrice *model.RetailPrice /*productID string, currency string, price decimal.Decimal*/) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	//retailPrice := model.NewRetailPrice(productID, currency, price)

	_, err := shopDb.Exec(`INSERT INTO retail_prices (product_id, currency, retail_price, date_updated)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (product_id, currency) DO UPDATE SET
		retail_price = $3,
		date_updated = $4`,
		retailPrice.ProductID,
		retailPrice.Currency,
		1,
		retailPrice.DateUpdated,
	)

	if err != nil {
		return fmt.Errorf("failed to insert retail price : <%w>", err)
	}

	return nil
}

/*


func GetRetailPrice(productID string, currency string) (*model.RetailPrice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	filter := bson.D{
		primitive.E{Key: "product_id", Value: productID},
		primitive.E{Key: "currency", Value: currency},
	}

	r := retailPriceCollection.FindOne(ctx, filter)

	price := model.RetailPrice{}
	if err := r.Decode(&price); err != nil {
		return nil, err
	}

	return &price, nil
}
*/

func GetRetailPrice(productID string, currency string) (*model.RetailPrice, error) {
	if shopDb == nil {
		return nil, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT retail_price, date_updated FROM retail_prices WHERE product_id = $1 AND currency = $2;`
	row := shopDb.QueryRow(query, productID, currency)

	var price string
	var dateUpdated time.Time

	err := row.Scan(&price, &dateUpdated)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row in GetRetailPrice: <%w>", err)
	}

	p, err := decimal.NewFromString(price)
	if err != nil {
		return nil, fmt.Errorf("failed to decode price in GetRetailPrice: <%w>", err)
	}

	retailPrice := model.NewRetailPrice(productID, currency, p)

	return retailPrice, nil
}
