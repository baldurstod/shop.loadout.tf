package shop

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"shop.loadout.tf/src/server/model"
)

const bcryptCost = 14

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
		return nil, errors.New("failed to create a user id")
	}

	user := model.NewUser()
	user.Username = username
	//user.Password = password
	user.DisplayName = username
	user.ID = id

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: <%w>", err)
	}

	if err = insertUser(user, hashedPassword); err != nil {
		return nil, fmt.Errorf("failed to create a user: <%w>", err)
	}

	return user, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func insertUser(user *model.User, password string) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	address, err := json.Marshal(&user.Address)
	if err != nil {
		return fmt.Errorf("failed to marshal user.Address: <%w>", err)
	}

	orders := make([]string, 0, len(user.Orders))
	for order := range user.Orders {
		orders = append(orders, order)
	}

	favorites := make([]string, 0, len(user.Favorites))
	for favorite := range user.Favorites {
		favorites = append(favorites, favorite)
	}

	_, err = shopDb.Exec(`INSERT INTO users (id, username, password, display_name, email_verified, address, currency, orders, favorites, cart, date_created, date_updated)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		user.ID,
		user.Username,
		password,
		user.DisplayName,
		user.EmailVerified,
		address,
		user.Currency,
		orders,
		favorites,
		"{}",
		user.DateCreated,
		user.DateUpdated,
	)

	if err != nil {
		return fmt.Errorf("failed to insert user: <%w>", err)
	}

	return nil
}

func FindUserByID(userId string) (*model.User, error) {
	query := `SELECT id, username, password, display_name, email_verified, address, currency, orders, favorites, cart, date_created, date_updated FROM users WHERE id = $1;`

	user, _, err := findUser(query, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UserIDExist(id string) (bool, error) {
	if shopDb == nil {
		return false, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT username FROM users WHERE id = $1;`
	row := shopDb.QueryRow(query, id)

	var username string
	err := row.Scan(&username)
	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func UsernameExist(username string) (bool, error) {
	if shopDb == nil {
		return false, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT id FROM users WHERE username = $1;`
	row := shopDb.QueryRow(query, username)

	var id string
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func FindUserByName(username string, password string) (*model.User, error) {
	query := `SELECT id, username, password, display_name, email_verified, address, currency, orders, favorites, cart, date_created, date_updated FROM users WHERE username = $1;`

	user, hashedPassword, err := findUser(query, username)
	if err != nil {
		return nil, err
	}

	if !checkPasswordHash(password, hashedPassword) {
		return nil, WrongPasswordError
	}

	return user, nil
}

func findUser(query string, args ...any) (*model.User, string, error) {
	if shopDb == nil {
		return nil, "", errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	row := shopDb.QueryRow(query, args...)

	var hashedPassword string
	var orders []string
	var favorites []string
	var cart string
	var address string

	user := model.NewUser()

	err := row.Scan(&user.ID, &user.Username, &hashedPassword, &user.DisplayName, &user.EmailVerified, &address, &user.Currency, pq.Array(&orders), pq.Array(&favorites), &cart, &user.DateCreated, &user.DateUpdated)
	if err != nil {
		return nil, "", err
	}

	for _, order := range orders {
		user.AddOrder(order)
	}

	for _, favorite := range favorites {
		user.AddFavorite(favorite)
	}

	if err = json.Unmarshal([]byte(cart), &user.Cart); err != nil {
		return nil, "", err
	}

	if err = json.Unmarshal([]byte(address), &user.Address); err != nil {
		return nil, "", err
	}

	return user, hashedPassword, nil
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
	if len(favorites) == 0 {
		return nil
	}
	user, err := FindUserByID(userID)
	if err != nil {
		return err
	}

	for favorite := range favorites {
		user.AddFavorite(favorite)
	}

	err = updateFavorites(user)
	if err != nil {
		return err
	}

	return nil
}

func updateFavorites(user *model.User) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	favorites := make([]string, 0, len(user.Favorites))
	for favorite := range user.Favorites {
		favorites = append(favorites, favorite)
	}
	user.DateUpdated = time.Now()

	query := `UPDATE users SET favorites = $2, date_updated = $3 WHERE id = $1;`
	_, err := shopDb.Exec(query, user.ID, favorites, user.DateUpdated)

	if err != nil {
		return fmt.Errorf("failed to update favorites:  <%w>", err)
	}

	return nil
}

func SetUserCart(userID string, cart model.Cart) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `UPDATE users SET cart = $2, date_updated = $3 WHERE id = $1;`
	_, err := shopDb.Exec(query, userID, cart, time.Now())

	if err != nil {
		return fmt.Errorf("failed to update user cart:  <%w>", err)
	}

	return nil
}

func ClearUserCart(userId string) error {
	return SetUserCart(userId, model.Cart{})
}

func SetUserCurrency(userID string, currency string) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `UPDATE users SET currency = $2, date_updated = $3 WHERE id = $1;`
	_, err := shopDb.Exec(query, userID, currency, time.Now())

	if err != nil {
		return fmt.Errorf("failed to update user currency:  <%w>", err)
	}

	return nil
}

type UpdateUserFields struct {
	DisplayName string
	AddOrder    string
}

func UserAddOrder(userId string, orderId string) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	_, err := shopDb.Exec(`UPDATE users SET orders = array_append(orders, $2), date_updated = $3 WHERE id = $1;`,
		userId,
		orderId,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to update user: <%w>", err)
	}

	return nil
}

func SetUserDisplayName(userId string, displayName string) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	_, err := shopDb.Exec(`UPDATE users SET display_name = $2, date_updated = $3 WHERE id = $1;`,
		userId,
		displayName,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to update user: <%w>", err)
	}

	return nil
}
