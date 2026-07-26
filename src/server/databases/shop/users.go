package shop

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
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

	addressEncryptedField, addressEncryptedKey, addressKekId, err := enveloped.EncryptAES(address)
	if err != nil {
		return fmt.Errorf("failed to encrypt address: <%w>", err)
	}

	orders := make([]string, 0, len(user.Orders))
	for order := range user.Orders {
		orders = append(orders, order)
	}

	favorites := make([]string, 0, len(user.Favorites))
	for favorite := range user.Favorites {
		favorites = append(favorites, favorite)
	}

	cartItems, err := json.Marshal(&user.Cart.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal user.Cart.Items: <%w>", err)
	}

	_, err = shopDb.Exec(`INSERT INTO users (id, username, password, display_name, email, email_verified, address, address_dek, address_kek, currency, orders, favorites, cart_items, date_created, date_updated)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
		user.ID,
		user.Username,
		password,
		user.DisplayName,
		user.Email,
		user.EmailVerified,
		addressEncryptedField,
		addressEncryptedKey,
		addressKekId,
		user.Currency,
		orders,
		pq.Array(favorites),
		cartItems,
		user.DateCreated,
		user.DateUpdated,
	)

	if err != nil {
		return fmt.Errorf("failed to insert user: <%w>", err)
	}

	return nil
}

func FindUserByID(userId string) (*model.User, error) {
	query := `SELECT id, username, password, display_name, email, email_verified, address, address_dek, address_kek, currency, orders, favorites, cart_items, date_created, date_updated FROM users WHERE id = $1;`

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
	query := `SELECT id, username, password, display_name, email, email_verified, address, address_dek, address_kek, currency, orders, favorites, cart_items, date_created, date_updated FROM users WHERE username = $1;`

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
	var cartItems string
	var encryptedAddress string
	var encryptedAddressDek string
	var addressKekId int64

	user := model.NewUser()

	err := row.Scan(&user.ID, &user.Username, &hashedPassword, &user.DisplayName, &user.Email, &user.EmailVerified, &encryptedAddress, &encryptedAddressDek, &addressKekId, &user.Currency, pq.Array(&orders), pq.Array(&favorites), &cartItems, &user.DateCreated, &user.DateUpdated)
	if err != nil {
		return nil, "", err
	}

	for _, order := range orders {
		user.AddOrder(order)
	}

	for _, favorite := range favorites {
		user.AddFavorite(favorite)
	}

	if err = json.Unmarshal([]byte(cartItems), &user.Cart.Items); err != nil {
		return nil, "", err
	}

	/*
		plainAddressDek, err := enveloped.DecryptDek(context.Background(), []byte(encryptedAddressDek))
		if err != nil {
			return nil, "", fmt.Errorf("failed to decrypt DEK: <%w>", err)
		}
	*/

	addressDecryptedField, err := enveloped.DecryptAES([]byte(encryptedAddress), []byte(encryptedAddressDek), addressKekId)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decrypt shipping address: <%w>", err)
	}

	if err = json.Unmarshal([]byte(addressDecryptedField), &user.Address); err != nil {
		return nil, "", err
	}

	return user, hashedPassword, nil
}

func SetUserFavorite(userID string, productID string, isFavorite bool) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	var query string
	if isFavorite {
		query = `UPDATE users SET favorites = array_append(favorites, $2), date_updated = $3 WHERE id = $1;`
	} else {
		query = `UPDATE users SET favorites = array_remove(favorites, $2), date_updated = $3 WHERE id = $1;`
	}

	res, err := shopDb.Exec(query, userID, productID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update user favorites:  <%w>", err)
	}

	if rows, err := res.RowsAffected(); rows != 1 || err != nil {
		if err != nil {
			return fmt.Errorf("failed to get rows affected %s: <%w>", userID, err)
		} else {
			return fmt.Errorf("failed to update user %s: %d rows affected, expected 1", userID, rows)
		}
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

	err = UpdateUser(*user, UpdateUserFields{Favorites: true})
	if err != nil {
		return err
	}

	return nil
}

func ClearUserCart(userId string) error {
	return UpdateUser(model.User{ID: userId}, UpdateUserFields{Cart: true})
}

func UserAddOrder(userId string, orderId string) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	res, err := shopDb.Exec(`UPDATE users SET orders = array_append(orders, $2), date_updated = $3 WHERE id = $1;`,
		userId,
		orderId,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to update user: <%w>", err)
	}

	if rows, err := res.RowsAffected(); rows != 1 || err != nil {
		if err != nil {
			return fmt.Errorf("failed to get rows affected %s: <%w>", userId, err)
		} else {
			return fmt.Errorf("failed to update user %s: %d rows affected, expected 1", userId, rows)
		}
	}

	return nil
}

type UpdateUserFields struct {
	Username      bool
	DisplayName   bool
	Email         bool
	EmailVerified bool
	Favorites     bool
	Currency      bool
	Cart          bool
	Address       bool
}

func UpdateUser(user model.User, fields UpdateUserFields) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	queryString := make([]string, 0)
	queryParams := []any{user.ID, time.Now()}

	v := reflect.ValueOf(fields)
	typeOfS := v.Type()

	// Using reflection to list UpdateUserFields fields
	for i := 0; i < v.NumField(); i++ {
		name := typeOfS.Field(i).Name
		value := v.Field(i).Bool()
		if !value {
			continue
		}

		addSetStatement := func(column string, value any) {
			param := "$" + strconv.Itoa(len(queryParams)+1)

			queryString = append(queryString, column+" = "+param)
			queryParams = append(queryParams, value)
		}

		switch name {
		case "Username":
			addSetStatement("username", user.Username)
		case "DisplayName":
			addSetStatement("display_name", user.DisplayName)
		case "Email":
			addSetStatement("email", user.Email)
		case "EmailVerified":
			addSetStatement("email_verified", user.EmailVerified)
		case "Currency":
			addSetStatement("currency", user.Currency)
		case "Favorites":

			favorites := make([]string, 0, len(user.Favorites))
			for favorite := range user.Favorites {
				favorites = append(favorites, favorite)
			}
			addSetStatement("favorites", pq.Array(favorites))
		case "Address":
			address, err := json.Marshal(&user.Address)
			if err != nil {
				return fmt.Errorf("failed to marshal user.Address: <%w>", err)
			}

			addressEncryptedField, addressEncryptedKey, addressKekId, err := enveloped.EncryptAES(address)
			if err != nil {
				return fmt.Errorf("failed to encrypt user address: <%w>", err)
			}

			addSetStatement("address", addressEncryptedField)
			addSetStatement("address_dek", addressEncryptedKey)
			addSetStatement("address_kek", addressKekId)
		case "Cart":

			cartItems, err := json.Marshal(&user.Cart.Items)
			if err != nil {
				return fmt.Errorf("failed to marshal user.Cart: <%w>", err)
			}

			addSetStatement("cart_items", cartItems)
		default:
			return errors.New("missing field in UpdateUser " + name)
		}
	}

	if len(queryString) == 0 {
		return errors.New("failed to update user: no field selected for update")
	}

	query := `UPDATE users SET date_updated = $2,` + strings.Join(queryString, ",") + ` WHERE id = $1;`
	res, err := shopDb.Exec(query, queryParams...)

	if err != nil {
		return fmt.Errorf("failed to update user: <%w>", err)
	}

	if rows, err := res.RowsAffected(); rows != 1 || err != nil {
		if err != nil {
			return fmt.Errorf("failed to get rows affected %s: <%w>", user.ID, err)
		} else {
			return fmt.Errorf("failed to update user %s: %d rows affected, expected 1", user.ID, rows)
		}
	}

	return nil
}
