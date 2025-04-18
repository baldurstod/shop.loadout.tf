package databases

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
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
var ordersCollection2 *mongo.Collection
var retailPriceCollection *mongo.Collection
var mockupTasksCollection *mongo.Collection
var usersCollection *mongo.Collection
var usersDecryptCollection *mongo.Collection

var secureClient *mongo.Client
var clientEnc *mongo.ClientEncryption
var dataKeyId primitive.Binary

const maxCreationAttempts = 10
const MongoTimeout = 30 * time.Second

func InitShopDB(config config.Database) {
	ctx, cancelConnect := context.WithCancel(context.Background())
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectURI))
	if err != nil {
		panic(err)
	}

	defer closeMongoDB(cancelConnect)

	productsCollection = client.Database(config.DBName).Collection("products")
	contactsCollection = client.Database(config.DBName).Collection("contacts")
	ordersCollection = client.Database(config.DBName).Collection("orders")
	retailPriceCollection = client.Database(config.DBName).Collection("retail_price")
	mockupTasksCollection = client.Database(config.DBName).Collection("mockup_tasks")
	usersCollection = client.Database(config.DBName).Collection("users")

	createUniqueIndex(productsCollection, "id", []string{"id"}, true)
	createUniqueIndex(ordersCollection, "id", []string{"id"}, true)
	createUniqueIndex(retailPriceCollection, "product_id,currency", []string{"product_id", "currency"}, true)
	createUniqueIndex(usersCollection, "id", []string{"id"}, true)
	createUniqueIndex(usersCollection, "email", []string{"email"}, true)

	if err := initEncryption(config); err != nil {
		panic(fmt.Errorf("error while initializing encryption %w", err))
	}

	ordersCollection2 = secureClient.Database(config.DBName).Collection("orders")
	usersDecryptCollection = secureClient.Database(config.DBName).Collection("users")
}

func initEncryption(config config.Database) error {
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(config.KeyVault.DEKID)))
	len, err := base64.StdEncoding.Decode(base64Text, []byte(config.KeyVault.DEKID))
	if err != nil {
		return fmt.Errorf("unable to decode key id: %s %w", config.KeyVault.DEKID, err)
	}
	dataKeyId = primitive.Binary{Subtype: 0x04, Data: base64Text[:len]}

	// Init KMS config
	provider := "kmip"
	kmsProviders := map[string]map[string]any{
		provider: {
			"endpoint": config.KeyVault.KMS.Endpoint,
		},
	}

	// Init TLS config
	tlsConfig := make(map[string]*tls.Config)
	tlsOpts := map[string]any{
		"tlsCertificateKeyFile": config.KeyVault.KMS.CertificatePath,
	}
	kmipConfig, err := options.BuildTLSConfig(tlsOpts)
	if err != nil {
		return err
	}
	tlsConfig["kmip"] = kmipConfig

	keyVaultNamespace := config.KeyVault.DBName + "." + config.KeyVault.Collection

	autoEncryptionOpts := options.AutoEncryption().SetKmsProviders(kmsProviders).SetKeyVaultNamespace(keyVaultNamespace).SetTLSConfig(tlsConfig).SetBypassAutoEncryption(true)

	secureClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(config.ConnectURI).SetAutoEncryptionOptions(autoEncryptionOpts))
	if err != nil {
		return fmt.Errorf("connect error for regular client: %w", err)
	}

	clientEncryptionOpts := options.ClientEncryption().SetKeyVaultNamespace(keyVaultNamespace).SetKmsProviders(kmsProviders).SetTLSConfig(tlsConfig)

	clientEnc, err = mongo.NewClientEncryption(secureClient, clientEncryptionOpts)
	if err != nil {
		return fmt.Errorf("newClientEncryption error %w", err)
	}

	return nil
}

func Cleanup() {
	if secureClient != nil {
		secureClient.Disconnect(context.TODO())
	}
	if clientEnc != nil {
		clientEnc.Close(context.TODO())
	}
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
		//log.Println("Failed to create index", name, "on collection", collection.Name(), err)
		panic(fmt.Errorf("failed to create index %s on collection %w", collection.Name(), err))
	}
}

func closeMongoDB(c context.CancelFunc) {
	if c != nil {
		c()
	}
}

func GetProduct(productID string) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	filter := bson.D{
		primitive.E{Key: "id", Value: productID},
		primitive.E{Key: "status", Value: "completed"},
	}

	cursor, err := productsCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("unable to find product %s: %w", productID, err)
	}

	for cursor.Next(context.TODO()) {
		product := model.NewProduct()
		if err := cursor.Decode(&product); err != nil {
			return nil, fmt.Errorf("error while decoding cursor for product %s: %w", productID, err)
		} else {
			return &product, nil
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor returned an error for product %s: %w", productID, err)
	}

	return nil, fmt.Errorf("product %s not found: %w", productID, err)
}

func GetProducts() ([]*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	filter := bson.D{primitive.E{Key: "status", Value: "completed"}}

	cursor, err := productsCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("unable to find products: %w", err)
	}

	results := []*model.Product{}
	variants := make(map[string]any)

	for cursor.Next(context.TODO()) {
		product := model.NewProduct()
		if err := cursor.Decode(&product); err != nil {
			return nil, fmt.Errorf("error while decoding cursor: %w", err)
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
		return nil, fmt.Errorf("cursor returned an error: %w", err)
	}

	return results, nil
}

func SendContact(subject string, email string, content string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	insertOneResult, err := contactsCollection.InsertOne(ctx, bson.M{
		"subject":      subject,
		"email":        email,
		"content":      content,
		"date_created": time.Now().Unix(),
		"status":       "created",
	})

	if err != nil {
		return "", err
	}

	return insertOneResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func CreateProduct() (*model.Product, error) {
	var id string
	ok := false
	for range maxCreationAttempts {
		id = createRandID()
		exist, err := ProductIDExist(id)
		if err != nil {
			return nil, err
		}

		if !exist {
			ok = true
			break
		}
	}

	if !ok {
		return nil, errors.New("unable to create an id")
	}

	product := model.NewProduct()
	product.ID = id

	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()
	if _, err := productsCollection.InsertOne(ctx, product); err != nil {
		return nil, err
	}

	return &product, nil
}

func UpdateProduct(product *model.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	filter := bson.D{primitive.E{Key: "id", Value: productID}}

	r := productsCollection.FindOne(ctx, filter)

	product := model.NewProduct()
	if err := r.Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

func createRandID() string {
	return randstr.String(12, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}
