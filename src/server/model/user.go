package model

type User struct {
	ID            string         `json:"id" bson:"id"`
	Username      string         `json:"username" bson:"username"`
	Password      string         `json:"password" bson:"password"`
	DateCreated   int64          `json:"date_created" bson:"date_created"`
	DateUpdated   int64          `json:"date_updated" bson:"date_updated"`
	EmailVerified bool           `json:"email_verified" bson:"email_verified"`
	Orders        map[string]any `json:"orders" bson:"orders"`
	Favorites     map[string]any `json:"favorites" bson:"favorites"`
	Currency      string         `json:"currency" bson:"currency"`
	Cart
	Address
}

func NewUser(username string, password string) *User {
	return &User{
		Username:      username,
		Password:      password,
		EmailVerified: false,
		Orders:        map[string]any{},
		Favorites:     map[string]any{},
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
