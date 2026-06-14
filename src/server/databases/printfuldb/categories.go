package printfuldb

import (
	"errors"
	"fmt"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
)

type MongoCategory struct {
	ID          int                    `json:"id" bson:"id"`
	LastUpdated int64                  `json:"last_updated" bson:"last_updated"`
	Category    printfulmodel.Category `json:"category" bson:"category"`
}

func GetCategories(language string) ([]printfulmodel.Category, error) {
	if printfulDb == nil {
		return nil, errors.New("database is not initialized. Did you forgot to call openPostgre ?")
	}

	query := `SELECT id, parent_id, image_url, title FROM categories WHERE language = $1;`
	res, err := printfulDb.Query(query, language)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query "+query+"in GetCategories: <%w>", err)
	}
	defer res.Close()

	categories := make([]printfulmodel.Category, 0, 400)
	for res.Next() {
		var id int
		var parent_id int
		var image_url string
		var title string

		err = res.Scan(&id, &parent_id, &image_url, &title)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row in GetCategories: <%w>", err)
		}
		doc := printfulmodel.Category{ID: id, ParentID: parent_id, ImageURL: image_url, Title: title}

		categories = append(categories, doc)
	}

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("failed to get next row in GetCategories: <%w>", err)
	}

	return categories, nil
}
