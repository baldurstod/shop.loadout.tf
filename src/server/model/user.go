package model

import (
	"time"

	"shop.loadout.tf/src/server/constants"
)

type User struct {
	ID            string              `json:"id" bson:"id"`
	Username      string              `json:"username" bson:"username"`
	Password      string              `json:"password" bson:"password"`
	DisplayName   string              `json:"display_name" bson:"display_name"`
	DateCreated   int64               `json:"date_created" bson:"date_created"`
	DateUpdated   int64               `json:"date_updated" bson:"date_updated"`
	EmailVerified bool                `json:"email_verified" bson:"email_verified"`
	Orders        map[string]struct{} `json:"orders" bson:"orders"`
	Favorites     map[string]struct{} `json:"favorites" bson:"favorites"`
	Currency      string              `json:"currency" bson:"currency"`
	Cart
	Address
}

func NewUser(username string, password string) *User {
	return &User{
		Username:      username,
		Password:      password,
		DisplayName:   username,
		DateCreated:   time.Now().Unix(),
		DateUpdated:   time.Now().Unix(),
		EmailVerified: false,
		Orders:        map[string]struct{}{},
		Favorites:     map[string]struct{}{},
		Currency:      constants.DEFAULT_CURRENCY,
		Cart:          NewCart(),
	}
}

func (user *User) AddOrder(id string) {
	user.Orders[id] = struct{}{}
}

func (user *User) AddFavorite(id string) {
	user.Favorites[id] = struct{}{}
}

func (user *User) RemoveFavorite(id string) {
	delete(user.Favorites, id)
}
