package mongo

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/baldurstod/randstr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/model"
)

var productsCollection *mongo.Collection
var contactsCollection *mongo.Collection
var ordersCollection *mongo.Collection
var retailPriceCollection *mongo.Collection
var mockupTasksCollection *mongo.Collection

func InitShopDB(config config.Database) {
	ctx, cancelConnect := context.WithCancel(context.Background())
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectURI))
	if err != nil {
		log.Println(err)
		panic(err)
	}

	defer closeMongoDB(cancelConnect)

	productsCollection = client.Database(config.DBName).Collection("products")
	contactsCollection = client.Database(config.DBName).Collection("contacts")
	ordersCollection = client.Database(config.DBName).Collection("orders")
	retailPriceCollection = client.Database(config.DBName).Collection("retail_price")
	mockupTasksCollection = client.Database(config.DBName).Collection("mockup_tasks")

	createUniqueIndex(productsCollection, "id", []string{"id"}, true)
	createUniqueIndex(ordersCollection, "id", []string{"id"}, true)
	createUniqueIndex(retailPriceCollection, "product_id,currency", []string{"product_id", "currency"}, true)

	if err := initEncryption(config.KeyVault); err != nil {
		log.Println(err)
		panic(err)
	}
}

func initEncryption(vault config.KeyVault) error {
	keyVaultClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(vault.ConnectURI))
	if err != nil {
		return fmt.Errorf("connect error for regular client: %v", err)
	}
	defer func() {
		keyVaultClient.Disconnect(context.TODO())
	}()

	keyVaultNamespace := vault.DBName + "." + vault.Collection

	// Init TLS config
	tlsConfig := make(map[string]*tls.Config)
	tlsOpts := map[string]interface{}{
		"tlsCertificateKeyFile": vault.KMS.CertificatePath,
	}
	kmipConfig, err := options.BuildTLSConfig(tlsOpts)
	if err != nil {
		return err
	}
	tlsConfig["kmip"] = kmipConfig

	// Init KMS config
	provider := "kmip"
	kmsProviders := map[string]map[string]any{
		provider: {
			"endpoint": vault.KMS.Endpoint,
		},
	}

	clientEncryptionOpts := options.ClientEncryption().SetKeyVaultNamespace(keyVaultNamespace).SetKmsProviders(kmsProviders).SetTLSConfig(tlsConfig)

	clientEnc, err := mongo.NewClientEncryption(keyVaultClient, clientEncryptionOpts)
	if err != nil {
		return fmt.Errorf("newClientEncryption error %v", err)
	}
	defer func() {
		clientEnc.Close(context.TODO())
	}()

	return nil
}

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
		log.Println("Failed to create index", name, "on collection", collection.Name(), err)
	}
}

func closeMongoDB(c context.CancelFunc) {
	if c != nil {
		c()
	}
}

func GetProduct(productID string) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		primitive.E{Key: "id", Value: productID},
		primitive.E{Key: "status", Value: "completed"},
	}

	cursor, err := productsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		product := model.NewProduct()
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		} else {
			return &product, nil
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return nil, errors.New("product not found")
}

func GetProducts() ([]*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{primitive.E{Key: "status", Value: "completed"}}

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

		if _, ok := variants[product.ID]; ok {
			continue
		}

		for _, variantID := range product.VariantIDs {
			variants[variantID] = struct{}{}
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

func CreateProduct() (*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	product := model.NewProduct()
	product.ID = randstr.String(12, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	_, err := productsCollection.InsertOne(ctx, product)
	if mongo.IsDuplicateKeyError(err) {
		return CreateProduct() // TODO: improve that
	}

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func UpdateProduct(product *model.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Replace().SetUpsert(true)
	product.DateUpdated = time.Now().Unix()

	filter := bson.D{primitive.E{Key: "id", Value: product.ID}}
	_, err := productsCollection.ReplaceOne(ctx, filter, product, opts)
	if err != nil {
		return err
	}

	return nil
}

func FindProduct(productID string) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{primitive.E{Key: "id", Value: productID}}

	r := productsCollection.FindOne(ctx, filter)

	product := model.NewProduct()
	if err := r.Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
}
