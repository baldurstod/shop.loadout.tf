package shop

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

func InsertPaypalError(content map[string]any) (int64, error) {
	if shopDb == nil {
		return 0, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	contentJson, err := json.Marshal(content)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal content in InsertPaypalError: <%w>", err)
	}

	var id int64
	err = shopDb.QueryRow(`INSERT INTO paypal_errors (content, date_created) VALUES ($1, $2) RETURNING id`,
		contentJson,
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to insert paypal error: <%w>", err)
	}

	return id, nil
}
