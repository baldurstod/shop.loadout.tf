package databases

import (
	"context"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"shop.loadout.tf/src/server/model"
)

func SetRetailPrice(productID string, currency string, price decimal.Decimal) (*model.RetailPrice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	retailPrice := model.NewRetailPrice(productID, currency, price)

	retailPrice.SetPrice(price)

	/*

		ProductID   primitive.ObjectID `json:"product_id" bson:"product_id"`
		Currency    string             `json:"currency" bson:"currency"`
		RetailPrice decimal.Decimal    `json:"retail_price" bson:"retail_price"`
	*/

	opts := options.Replace().SetUpsert(true)

	filter := bson.D{primitive.E{Key: "product_id", Value: productID}, primitive.E{Key: "currency", Value: currency}}
	_, err := retailPriceCollection.ReplaceOne(ctx, filter, retailPrice, opts)
	if err != nil {
		return nil, err
	}

	return retailPrice, nil
}

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
