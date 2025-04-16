package databases

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"shop.loadout.tf/src/server/model"
)

func CreateUser(email string, password string) (*model.User, error) {
	emailExist, err := UserEmailExist(email)
	if err != nil {
		return nil, err
	}

	if emailExist {
		return nil, fmt.Errorf("email %s already exist", email)
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
		return nil, errors.New("unable to create an id")
	}

	user := model.NewUser(email, password)
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

func UserEmailExist(email string) (bool, error) {
	r := usersDecryptCollection.FindOne(context.Background(), bson.D{primitive.E{Key: "address.email", Value: email}})

	err := r.Err()

	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func FindUserByEmail(email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	filter := bson.D{primitive.E{Key: "address.email", Value: email}}

	r := usersDecryptCollection.FindOne(ctx, filter)

	user := model.User{}
	if err := r.Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
