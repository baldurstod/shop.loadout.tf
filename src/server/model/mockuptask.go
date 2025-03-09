package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
)

type MockupTask struct {
	ID          primitive.ObjectID             `json:"_id,omitempty" bson:"_id,omitempty"`
	ProductIDs  []string                       `json:"product_ids" bson:"product_ids"`
	SourceImage string                         `json:"source_image,omitempty" bson:"source_image,omitempty"`
	Template    *printfulmodel.MockupTemplates `json:"template,omitempty" bson:"template,omitempty"`
	Status      string                         `json:"status" bson:"status"`
	DateCreated int64                          `json:"date_created" bson:"date_created"`
	DateUpdated int64                          `json:"date_updated" bson:"date_updated"`
}

func NewMockupTask() MockupTask {
	return MockupTask{
		ProductIDs: []string{},
	}
}

func (task *MockupTask) AddProduct(productID string) {
	task.ProductIDs = append(task.ProductIDs, productID)
}
