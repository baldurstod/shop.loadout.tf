package printfuldb

import (
	"context"
	"time"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"go.mongodb.org/mongo-driver/bson"
)

func FindProductsPrices(currency string) ([]*printfulmodel.ProductPrices, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
