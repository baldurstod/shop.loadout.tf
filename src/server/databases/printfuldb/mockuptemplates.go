package printfuldb

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
)

func FindMockupTemplates(productID int) ([]printfulmodel.MockupTemplates, bool, error) {
	if printfulDb == nil {
		return nil, false, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT mockup_templates, last_updated FROM mockup_templates WHERE product_id = $1;`
	row := printfulDb.QueryRow(query, productID)

	var mockupTemplates string
	var lastUpdated int64

	err := row.Scan(&mockupTemplates, &lastUpdated)
	if err != nil {
		return nil, false, fmt.Errorf("failed to scan row in FindMockupTemplates: <%w>", err)
	}

	templates := []printfulmodel.MockupTemplates{}
	if err = json.Unmarshal([]byte(mockupTemplates), &templates); err != nil {
		return nil, false, err
	}

	return templates, time.Now().Unix()-lastUpdated > cacheMaxAge, nil
}
