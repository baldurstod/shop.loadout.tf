package shop

import (
	"errors"
	"time"

	"github.com/baldurstod/randstr"
	"go.mongodb.org/mongo-driver/mongo"
	"shop.loadout.tf/src/server/model"
)

var usersCollection *mongo.Collection
var taxCollection *mongo.Collection

const maxCreationAttempts = 10
const MongoTimeout = 30 * time.Second

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
	if err := insertProduct(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

func UpdateProduct(product *model.Product) error {
	if err := insertProduct(product); err != nil {
		return err
	}
	return nil
}

func createRandID() string {
	return randstr.String(12, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}
