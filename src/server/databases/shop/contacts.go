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

	res, err := shopDb.Exec(`INSERT INTO contacts (subject, email, content, status, date_created)
						VALUES ($2, $3, $4, $5, $6)`,
		subject,
		email,
		content,
		"created",
		time.Now(),
	)

	if err != nil {
		return -1, fmt.Errorf("failed to insert contact: <%w>", err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return -1, fmt.Errorf("failed to get last inserted id: <%w>", err)
	}

	return id, nil
}
