package mongo

import (
	"context"
	_ "go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"shop.loadout.tf/src/server/config"
	_ "time"
)

var cancelConnect context.CancelFunc
var shopCollection *mongo.Collection

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
}

func closeMongoDB() {
	if cancelConnect != nil {
		cancelConnect()
	}
}
