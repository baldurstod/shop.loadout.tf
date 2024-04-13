package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
var contactsCollection *mongo.Collection
var ordersCollection *mongo.Collection

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
	contactsCollection = client.Database(config.DBName).Collection("contacts")
	ordersCollection = client.Database(config.DBName).Collection("orders")
}

func closeMongoDB() {
	if cancelConnect != nil {
		cancelConnect()
	}
}

func GetProduct(productID string) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	docID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{"_id", docID},
		{"status", "completed"},
	}

	cursor, err := productsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		product := model.NewProduct()
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}

		return &product, nil
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func GetProducts() ([]*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{"status", "completed"}}

	cursor, err := productsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	results := []*model.Product{}
	variants := make(map[string]interface{})

	for cursor.Next(context.TODO()) {
		product := model.NewProduct()
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}

		if _, ok := variants[product.ID.Hex()]; ok {
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

func SendContact(params map[string]interface{}) (string, error) {
	log.Println(params)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	insertOneResult, err := contactsCollection.InsertOne(ctx, bson.M{
		"subject":      params["subject"],
		"email":        params["email"],
		"content":      params["content"],
		"date_created": time.Now().Unix(),
		"status":       "created",
	})

	if err != nil {
		return "", err
	}

	return insertOneResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func CreateOrder() (*model.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	order := model.NewOrder()

	insertOneResult, err := ordersCollection.InsertOne(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	order.ID = insertOneResult.InsertedID.(primitive.ObjectID)

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

	filter := bson.D{{"_id", order.ID}}
	_, err := ordersCollection.ReplaceOne(ctx, filter, order, opts)
	if err != nil {
		return err
	}

	return nil
}

func CreateProduct() (*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	product := model.NewProduct()

	insertOneResult, err := productsCollection.InsertOne(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	product.ID = insertOneResult.InsertedID.(primitive.ObjectID)

	return &product, nil
}

func UpdateProduct(product *model.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Replace().SetUpsert(true)
	product.DateUpdated = time.Now().Unix()

	filter := bson.D{{"_id", product.ID}}
	_, err := productsCollection.ReplaceOne(ctx, filter, product, opts)
	if err != nil {
		return err
	}

	return nil
}

func FindProduct(productID string) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	docID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", docID}}

	r := productsCollection.FindOne(ctx, filter)

	product := model.NewProduct()
	if err := r.Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
}
