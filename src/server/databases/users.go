package databases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"shop.loadout.tf/src/server/model"
)

func CreateUser(username string, password string) (*model.User, error) {
	userExist, err := UsernameExist(username)
	if err != nil {
		return nil, err
	}

	if userExist {
		return nil, fmt.Errorf("username %s already exist", username)
	}

	var id string
	ok := false
	for range maxCreationAttempts {
		id = createRandID()
		exist, err := UserIDExist(id)
		if err != nil {
			return nil, err
		}

		if !exist {
			ok = true
			break
		}
	}

	if !ok {
		return nil, errors.New("unable to create a user id")
	}

	user := model.NewUser(username, password)
	user.ID = id

	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()
	if _, err = usersCollection.InsertOne(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func FindUserByID(id string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	filter := bson.D{primitive.E{Key: "id", Value: id}}

	r := usersDecryptCollection.FindOne(ctx, filter)

	user := model.User{}
	if err := r.Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func UserIDExist(id string) (bool, error) {
	r := usersDecryptCollection.FindOne(context.Background(), bson.D{primitive.E{Key: "id", Value: id}})

	err := r.Err()

	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func UsernameExist(username string) (bool, error) {
	r := usersDecryptCollection.FindOne(context.Background(), bson.D{primitive.E{Key: "username", Value: username}})

	err := r.Err()

	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func FindUserByName(username string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	filter := bson.D{primitive.E{Key: "username", Value: username}}

	r := usersDecryptCollection.FindOne(ctx, filter)

	user := model.User{}
	if err := r.Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(user *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	opts := options.Replace().SetUpsert(true)
	user.DateUpdated = time.Now().Unix()

	var err error
	/*
		encryptedUser, err := encryptUser(user)
		if err != nil {
			return err
		}
	*/

	filter := bson.D{primitive.E{Key: "id", Value: user.ID}}
	_, err = usersCollection.ReplaceOne(ctx, filter, user, opts)
	if err != nil {
		return err
	}

	return nil
}
func encryptUser(user *model.User) (*bson.M, error) {
	/*
		shippingAddressEncryptedField, err := encryptAddress(&user.ShippingAddress)
		if err != nil {
			return nil, err
		}

		billingAddressEncryptedField, err := encryptAddress(&user.BillingAddress)
		if err != nil {
			return nil, err
		}
	*/

	return &bson.M{
		"id":           user.ID,
		"currency":     user.Currency,
		"date_created": user.DateCreated,
		"date_updated": user.DateUpdated,
	}, nil
}
