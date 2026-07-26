package shop

import (
	"errors"
	"fmt"
	"time"
)

func InsertContact(subject string, email string, content string) (int64, error) {
	if shopDb == nil {
		return -1, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	var id int64
	err := shopDb.QueryRow(`INSERT INTO contacts (subject, email, content, status, date_created)
						VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		subject,
		email,
		content,
		"created",
		time.Now(),
	).Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("failed to insert contact: <%w>", err)
	}

	return id, nil
}
