package databases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		return nil, fmt.Errorf("failed to decode user: %w", err)
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

func SetUserFavorite(userID string, productID string, isFavorite bool) error {
	user, err := FindUserByID(userID)
	if err != nil {
		return err
	}

	if isFavorite {
		user.Favorites[productID] = struct{}{}
	} else {
		delete(user.Favorites, productID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	filter := bson.D{{Key: "id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "favorites", Value: user.Favorites}, {Key: "date_updated", Value: time.Now().Unix()}}}}
	_, err = usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func AddUserFavorites(userID string, favorites map[string]any) error {
	user, err := FindUserByID(userID)
	if err != nil {
		return err
	}

	for favorite := range favorites {
		user.AddFavorite(favorite)
	}

	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	filter := bson.D{{Key: "id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "favorites", Value: user.Favorites}, {Key: "date_updated", Value: time.Now().Unix()}}}}
	_, err = usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func SetUserCart(userID string, cart model.Cart) error {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	filter := bson.D{{Key: "id", Value: userID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "cart", Value: cart}, {Key: "date_updated", Value: time.Now().Unix()}}}}
	_, err := usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func ClearUserCart(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	filter := bson.D{{Key: "id", Value: userID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "cart.items", Value: map[string]uint{}}, {Key: "date_updated", Value: time.Now().Unix()}}}}
	_, err := usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func SetUserCurrency(userID string, currency string) error {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	filter := bson.D{{Key: "id", Value: userID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "currency", Value: currency}, {Key: "date_updated", Value: time.Now().Unix()}}}}
	_, err := usersCollection.UpdateOne(ctx, filter, update)
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

type UpdateUserFields struct {
	DisplayName string
	AddOrder    string
}

func UpdateUser(userID string, fields UpdateUserFields) error {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	filter := bson.M{"id": userID}
	updateFields := bson.M{}
	update := bson.M{"$set": updateFields}

	if fields.DisplayName != "" {
		updateFields["display_name"] = fields.DisplayName
	}

	if fields.AddOrder != "" {
		updateFields["orders."+fields.AddOrder] = bson.M{}
	}

	_, err := usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
