package mongo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"shop.loadout.tf/src/server/model"
)

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
