package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	_"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"shop.loadout.tf/src/server/config"
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

func GetProducts() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{"status", "completed"}}
	opts := options.Find().SetProjection(bson.M{"_id": 0})

	cursor, err := productsCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	/*for _, result := range results {
		log.Println(result)
	}*/

	//objectID := res.InsertedID.(primitive.ObjectID)
	return results, nil
}
