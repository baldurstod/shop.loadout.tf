package model

type Option struct {
	Name  string `json:"name" bson:"name"`
	Type  string `json:"type" bson:"type"`
	Value string `json:"value" bson:"value"`
}
