package databases

import (
	"context"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"shop.loadout.tf/src/server/model"
)

func SetTaxRate(countryCode string, stateCode string, postalCode string, city string, rate decimal.Decimal) (*model.TaxRate, error) {

	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	taxRate := model.NewTaxRate(countryCode, stateCode, postalCode, city, rate)

	opts := options.Replace().SetUpsert(true)

	filter := bson.D{
		primitive.E{Key: "country_code", Value: countryCode},
		primitive.E{Key: "state_code", Value: stateCode},
		primitive.E{Key: "postal_code", Value: postalCode},
		primitive.E{Key: "city", Value: city},
	}
	_, err := taxCollection.ReplaceOne(ctx, filter, taxRate, opts)
	if err != nil {
		return nil, err
	}

	return taxRate, nil
}

func GetTaxRate(countryCode string, stateCode string, postalCode string, city string) (decimal.Decimal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "country_code", Value: countryCode}},
			bson.D{{Key: "country_code", Value: ""}}},
		},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "state_code", Value: stateCode}},
			bson.D{{Key: "state_code", Value: ""}}},
		},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "postal_code", Value: postalCode}},
			bson.D{{Key: "postal_code", Value: ""}}},
		},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "city", Value: city}},
			bson.D{{Key: "city", Value: ""}}},
		},
	}

	r := taxCollection.FindOne(ctx, filter)

	taxRate := model.TaxRate{}
	if err := r.Decode(&taxRate); err != nil {
		return decimal.Decimal{}, err
	}

	return decimal.NewFromString(taxRate.Rate.String())
}
