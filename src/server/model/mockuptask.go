package model

import (
	"time"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
)

type MockupTask struct {
	ID          int64                          `json:"id,omitempty" bson:"id,omitempty"`
	ProductIDs  []string                       `json:"product_ids" bson:"product_ids"`
	SourceImage string                         `json:"source_image,omitempty" bson:"source_image,omitempty"`
	Template    *printfulmodel.MockupTemplates `json:"template,omitempty" bson:"template,omitempty"`
	Status      string                         `json:"status" bson:"status"`
	DateCreated time.Time                      `json:"date_created" bson:"date_created"`
	DateUpdated time.Time                      `json:"date_updated" bson:"date_updated"`
}

func NewMockupTask() MockupTask {
	return MockupTask{
		ProductIDs: []string{},
	}
}

func (task *MockupTask) AddProduct(productID string) {
	task.ProductIDs = append(task.ProductIDs, productID)
}
