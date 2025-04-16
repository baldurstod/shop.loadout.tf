package model

type User struct {
	ID       string `json:"id" bson:"id"`
	Password string `json:"password" bson:"password"`
	Address
}

func NewUser(email string, password string) *User {
	return &User{Password: password, Address: Address{Email: email}}
}
