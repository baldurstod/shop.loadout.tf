package model

type User struct {
	ID            string          `json:"id" bson:"id"`
	Username      string          `json:"username" bson:"username"`
	Password      string          `json:"password" bson:"password"`
	DateCreated   int64           `json:"date_created" bson:"date_created"`
	DateUpdated   int64           `json:"date_updated" bson:"date_updated"`
	EmailVerified bool            `json:"email_verified" bson:"email_verified"`
	Orders        map[string]bool `json:"orders" bson:"orders"`
	Favorites     map[string]bool `json:"favorites" bson:"favorites"`
	Address
}

func NewUser(username string, password string) *User {
	return &User{
		Password:      password,
		Username:      username,
		EmailVerified: false,
	}
}

func (user *User) AddOrder(id string) {
	user.Orders[id] = true
}

func (user *User) AddFavorite(id string) {
	user.Favorites[id] = true
}

func (user *User) RemoveFavorite(id string) {
	delete(user.Favorites, id)
}
