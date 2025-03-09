package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"shop.loadout.tf/src/server/model"
)

func SetRetailPrice(productID string, currency string, price decimal.Decimal) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	retailPrice := model.RetailPrice{
		ProductID: productID,
		Currency:  currency,
	}

	retailPrice.SetRetailPrice(price)

	/*

		ProductID   primitive.ObjectID `json:"product_id" bson:"product_id"`
		Currency    string             `json:"currency" bson:"currency"`
		RetailPrice decimal.Decimal    `json:"retail_price" bson:"retail_price"`
	*/

	opts := options.Replace().SetUpsert(true)
	retailPrice.DateUpdated = time.Now().Unix()

	filter := bson.D{primitive.E{Key: "ProductID", Value: productID}, primitive.E{Key: "Currency", Value: currency}}
	_, err := retailPriceCollection.ReplaceOne(ctx, filter, retailPrice, opts)
	if err != nil {
		return err
	}

	return nil
}

func GetRetailPrice(productID string, currency string) (*model.RetailPrice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	docID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		primitive.E{Key: "_id", Value: docID},
		primitive.E{Key: "currency", Value: currency},
	}

	cursor, err := retailPriceCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		product := model.RetailPrice{}
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		} else {
			return &product, nil
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return nil, errors.New("product price not found")
}
