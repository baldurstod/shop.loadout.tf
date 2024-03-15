package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	_"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/model"
	"time"
)

var cancelConnect context.CancelFunc
var shopCollection *mongo.Collection
var productsCollection *mongo.Collection

func InitMongoDB(config config.Database) {
	var ctx context.Context
	ctx, cancelConnect = context.WithCancel(context.Background())
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectURI))
	if err != nil {
		log.Println(err)
		panic(err)
	}

	defer closeMongoDB()

	shopCollection = client.Database(config.DBName).Collection("shop")
	productsCollection = client.Database(config.DBName).Collection("products")
}

func closeMongoDB() {
	if cancelConnect != nil {
		cancelConnect()
	}
}

func GetProducts() ([]*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{"status", "completed"}}

	cursor, err := productsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var results []*model.Product
	variants := make(map[string]interface{})

	for cursor.Next(context.TODO()) {
		product := model.Product{}
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}

		if _, ok := variants[product.Id]; ok {
			continue
		}

		for _, variantId := range product.VariantIds {
			variants[variantId] = struct{}{}
		}

		results = append(results, &product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
