package printfuldb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createUniqueIndex(collection *mongo.Collection, name string, keys []string, unique bool) {
	keysDoc := bson.D{}
	for _, key := range keys {
		keysDoc = append(keysDoc, bson.E{Key: key, Value: 1})
	}

	if _, err := collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    keysDoc,
			Options: options.Index().SetUnique(unique).SetName(name),
		},
	); err != nil {
		err := fmt.Errorf("failed to create index %s on collection %s: %w", name, collection.Name(), err)
		log.Println(err)
		panic(err)
	}
}
