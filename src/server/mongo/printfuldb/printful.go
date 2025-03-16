package printfuldb

import (
	"context"
	"fmt"
	"log"
	"time"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"shop.loadout.tf/src/server/config"
)

var pfProductsCollection *mongo.Collection
var pfProductsPricesCollection *mongo.Collection
var pfMockupTemplatesCollection *mongo.Collection
var pfMockupStylesCollection *mongo.Collection
var pfVariantsCollection *mongo.Collection
var pfCountriesCollection *mongo.Collection
var pfCategoriesCollection *mongo.Collection

var cacheMaxAge int64 = 86400

func InitPrintfulDB(config config.Database) {
	ctx, cancelPrintful := context.WithCancel(context.Background())
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectURI))
	if err != nil {
		log.Println(err)
		panic(err)
	}

	defer closePrintfulDB(cancelPrintful)

	pfProductsCollection = client.Database(config.DBName).Collection("products")
	pfProductsPricesCollection = client.Database(config.DBName).Collection("products_prices")
	pfMockupTemplatesCollection = client.Database(config.DBName).Collection("mockup_templates")
	pfMockupStylesCollection = client.Database(config.DBName).Collection("mockup_styles")
	pfVariantsCollection = client.Database(config.DBName).Collection("variants")
	pfCountriesCollection = client.Database(config.DBName).Collection("countries")
	pfCategoriesCollection = client.Database(config.DBName).Collection("categories")

	createUniqueIndex(pfProductsCollection, "id", []string{"id"}, true)
	createUniqueIndex(pfVariantsCollection, "id", []string{"id"}, true)
	createUniqueIndex(pfVariantsCollection, "variant.catalog_product_id", []string{"variant.catalog_product_id"}, false)
	createUniqueIndex(pfProductsPricesCollection, "product_id", []string{"product_id"}, false)
	createUniqueIndex(pfProductsPricesCollection, "currency", []string{"currency"}, false)
	createUniqueIndex(pfProductsPricesCollection, "product_id,currency", []string{"product_id", "currency"}, true)
	createUniqueIndex(pfMockupTemplatesCollection, "product_id", []string{"product_id"}, false)
	createUniqueIndex(pfMockupStylesCollection, "product_id", []string{"product_id"}, false)
	createUniqueIndex(pfCountriesCollection, "code", []string{"code"}, false)
	createUniqueIndex(pfCategoriesCollection, "id", []string{"id"}, true)
}

func closePrintfulDB(c context.CancelFunc) {
	if c != nil {
		c()
	}
}

type MongoProduct struct {
	ID          int                   `json:"id" bson:"id"`
	LastUpdated int64                 `json:"last_updated" bson:"last_updated"`
	Product     printfulmodel.Product `json:"product" bson:"product"`
}

func FindProducts() ([]printfulmodel.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{}

	cursor, err := pfProductsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	products := make([]printfulmodel.Product, 0, 400)
	for cursor.Next(context.TODO()) {
		doc := MongoProduct{}
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}

		products = append(products, doc.Product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func FindProduct(productID int) (*printfulmodel.Product, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "id", Value: productID}}

	r := pfProductsCollection.FindOne(ctx, filter)

	doc := MongoProduct{}
	if err := r.Decode(&doc); err != nil {
		return nil, false, err
	}

	return &doc.Product, time.Now().Unix()-doc.LastUpdated > cacheMaxAge, nil
}

func FindVariants(productID int) (variants []printfulmodel.Variant, outdated bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "variant.catalog_product_id", Value: productID}}
	outdated = false

	cursor, err := pfVariantsCollection.Find(ctx, filter)
	if err != nil {
		return nil, false, err
	}

	variants = make([]printfulmodel.Variant, 0, 20)
	for cursor.Next(context.TODO()) {
		doc := MongoVariant{}
		if err := cursor.Decode(&doc); err != nil {
			return nil, false, err
		}

		if time.Now().Unix()-doc.LastUpdated > cacheMaxAge {
			outdated = true
		}

		variants = append(variants, doc.Variant)
	}

	if err := cursor.Err(); err != nil {
		return nil, false, err
	}

	return variants, outdated, nil
}

func InsertProduct(product *printfulmodel.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	variantIds, err := getProductVariantIds(product.ID)
	if err != nil {
		return fmt.Errorf("error in InsertProduct: %w", err)
	}

	product.CatalogVariantIDs = variantIds
	product.VariantCount = len(variantIds)

	opts := options.Replace().SetUpsert(true)

	filter := bson.D{{Key: "id", Value: product.ID}}
	doc := MongoProduct{ID: product.ID, LastUpdated: time.Now().Unix(), Product: *product}
	_, err = pfProductsCollection.ReplaceOne(ctx, filter, doc, opts)

	return err
}

func getProductVariantIds(productId int) ([]int, error) {
	//Update product variants
	variantIDs := make([]int, 0, 20)
	variants, _, err := FindVariants(productId)
	if err != nil {
		return nil, fmt.Errorf("error while finding variants in getProductVariantIds: %w", err)
	}

	for _, variant := range variants {
		variantIDs = append(variantIDs, variant.ID)
	}

	return variantIDs, err
}

func UpdateProductVariantIds(id int, variantIds []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5005*time.Second)
	defer cancel()

	filter := bson.D{{Key: "id", Value: id}}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "product.catalog_variant_ids", Value: variantIds},
		{Key: "product.variant_count", Value: len(variantIds)},
	}}}

	//doc := MongoProduct{ID: product.ID, LastUpdated: time.Now().Unix(), Product: *product}
	_, err := pfProductsCollection.UpdateOne(ctx, filter, update)

	return err
}

type MongoProductPrices struct {
	ProductID     int                         `json:"product_id" bson:"product_id"`
	Currency      string                      `json:"currency" bson:"currency"`
	LastUpdated   int64                       `json:"last_updated" bson:"last_updated"`
	ProductPrices printfulmodel.ProductPrices `json:"product_prices" bson:"product_prices"`
}

func InsertProductPrices(productPrices *printfulmodel.ProductPrices) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Replace().SetUpsert(true)

	//	filter := bson.D{{Key: "id", Value: productPrices.Product.ID}}
	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "product_id", Value: productPrices.Product.ID}},
				bson.D{{Key: "currency", Value: productPrices.Currency}},
			},
		},
	}

	doc := MongoProductPrices{ProductID: productPrices.Product.ID, Currency: productPrices.Currency, LastUpdated: time.Now().Unix(), ProductPrices: *productPrices}
	_, err := pfProductsPricesCollection.ReplaceOne(ctx, filter, doc, opts)

	return err
}
func FindProductPrices(productID int, currency string) (*printfulmodel.ProductPrices, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "product_id", Value: productID}},
				bson.D{{Key: "currency", Value: currency}},
			},
		},
	}

	r := pfProductsPricesCollection.FindOne(ctx, filter)

	doc := MongoProductPrices{}
	if err := r.Decode(&doc); err != nil {
		return nil, false, err
	}

	return &doc.ProductPrices, time.Now().Unix()-doc.LastUpdated > cacheMaxAge, nil
}

type MongoMockupTemplates struct {
	ProductID       int                             `json:"product_id" bson:"product_id"`
	LastUpdated     int64                           `json:"last_updated" bson:"last_updated"`
	MockupTemplates []printfulmodel.MockupTemplates `json:"mockup_templates" bson:"mockup_templates"`
}

func InsertMockupTemplates(productID int, mockupTemplates []printfulmodel.MockupTemplates) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Replace().SetUpsert(true)

	filter := bson.D{{Key: "product_id", Value: productID}}

	doc := MongoMockupTemplates{ProductID: productID, LastUpdated: time.Now().Unix(), MockupTemplates: mockupTemplates}
	_, err := pfMockupTemplatesCollection.ReplaceOne(ctx, filter, doc, opts)

	return err
}

func FindMockupTemplates(productID int) ([]printfulmodel.MockupTemplates, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "product_id", Value: productID}}

	r := pfMockupTemplatesCollection.FindOne(ctx, filter)

	doc := MongoMockupTemplates{}
	if err := r.Decode(&doc); err != nil {
		return nil, false, err
	}

	return doc.MockupTemplates, time.Now().Unix()-doc.LastUpdated > cacheMaxAge, nil
}

type MongoMockupStyles struct {
	ProductID    int                          `json:"product_id" bson:"product_id"`
	LastUpdated  int64                        `json:"last_updated" bson:"last_updated"`
	MockupStyles []printfulmodel.MockupStyles `json:"mockup_styles" bson:"mockup_styles"`
}

func InsertMockupStyles(productID int, mockupStyles []printfulmodel.MockupStyles) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Replace().SetUpsert(true)

	filter := bson.D{{Key: "product_id", Value: productID}}

	doc := MongoMockupStyles{ProductID: productID, LastUpdated: time.Now().Unix(), MockupStyles: mockupStyles}
	_, err := pfMockupStylesCollection.ReplaceOne(ctx, filter, doc, opts)

	return err
}

func FindMockupStyles(productID int) ([]printfulmodel.MockupStyles, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "product_id", Value: productID}}

	r := pfMockupStylesCollection.FindOne(ctx, filter)

	doc := MongoMockupStyles{}
	if err := r.Decode(&doc); err != nil {
		return nil, false, err
	}

	return doc.MockupStyles, time.Now().Unix()-doc.LastUpdated > cacheMaxAge, nil
}

type MongoVariant struct {
	ID          int                   `json:"id" bson:"id"`
	LastUpdated int64                 `json:"last_updated" bson:"last_updated"`
	Variant     printfulmodel.Variant `json:"variant" bson:"variant"`
}

func InsertVariant(variant *printfulmodel.Variant) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Replace().SetUpsert(true)

	filter := bson.D{{Key: "id", Value: variant.ID}}
	doc := MongoVariant{ID: variant.ID, LastUpdated: time.Now().Unix(), Variant: *variant}
	_, err := pfVariantsCollection.ReplaceOne(ctx, filter, doc, opts)

	return err
}
