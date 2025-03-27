package mongo

import (
	"context"
	"time"

	"github.com/baldurstod/randstr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"shop.loadout.tf/src/server/model"
)

func CreateOrder() (*model.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	order := model.NewOrder()
	order.ID = randstr.String(12, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	_, err := ordersCollection.InsertOne(ctx, order)
	if mongo.IsDuplicateKeyError(err) {
		return CreateOrder() // TODO: improve that
	}

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func UpdateOrder(order *model.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	/*docID, err := primitive.ObjectIDFromHex(order.ID)
	if err != nil {
		return nil, err
	}*/

	opts := options.Replace().SetUpsert(true)
	order.DateUpdated = time.Now().Unix()

	filter := bson.D{primitive.E{Key: "id", Value: order.ID}}
	_, err := ordersCollection.ReplaceOne(ctx, filter, order, opts)
	if err != nil {
		return err
	}

	return nil
}

func FindOrder(orderID string) (*model.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{primitive.E{Key: "id", Value: orderID}}

	r := ordersCollection.FindOne(ctx, filter)

	order := model.Order{}
	if err := r.Decode(&order); err != nil {
		return nil, err
	}

	return &order, nil
}

func FindOrderByPaypalID(paypalID string) (*model.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{primitive.E{Key: "paypal_order_id", Value: paypalID}}

	r := ordersCollection.FindOne(ctx, filter)

	order := model.Order{}
	if err := r.Decode(&order); err != nil {
		return nil, err
	}

	return &order, nil
}
