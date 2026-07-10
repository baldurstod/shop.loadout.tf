package printfuldb

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
)

func FindMockupStyles(productID int) ([]printfulmodel.MockupStyles, bool, error) {

	if printfulDb == nil {
		return nil, false, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT mockup_styles, last_updated FROM mockup_styles WHERE product_id = $1;`
	row := printfulDb.QueryRow(query, productID)

	var mockupStyles string
	var lastUpdated int64

	err := row.Scan(&mockupStyles, &lastUpdated)
	if err != nil {
		return nil, false, fmt.Errorf("failed to scan row in FindMockupStyles: <%w>", err)
	}

	styles := []printfulmodel.MockupStyles{}
	if err = json.Unmarshal([]byte(mockupStyles), &styles); err != nil {
		return nil, false, err
	}

	return styles, time.Now().Unix()-lastUpdated > cacheMaxAge, nil
}
