package model

type User struct {
	ID      string `json:"id" bson:"id"`
	Address `mapstructure:",squash"`
}

func NewUser(email string) *User {
	return &User{Address: Address{Email: email}}
}
