package shop

import (
	"errors"
	"fmt"
	"time"

	shopErrors "shop.loadout.tf/src/server/errors"
)

const codeValidity = 15 * time.Minute

const ErrCodeExpired = shopErrors.ErrString("code expired")

func InsertEmailVerification(userId string, email string, code string) (err error) {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	_, err = shopDb.Exec(`INSERT INTO email_verification (user_id, email, code, date_created) VALUES ($1, $2, $3, $4)
	ON CONFLICT (user_id, email) DO UPDATE SET
			code = $3,
			date_created = $4`,
		userId,
		email,
		code,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to insert email verification: <%w>", err)
	}

	return nil
}

func CheckEmailVerification(code string, email string) (userId string, err error) {
	if shopDb == nil {
		return "", errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT user_id, date_created FROM email_verification WHERE code = $1 AND email = $2;`

	row := shopDb.QueryRow(query, code, email)

	var dateCreated time.Time

	if err = row.Scan(&userId, &dateCreated); err != nil {
		return "", err
	}

	if time.Since(dateCreated) > codeValidity {
		return "", ErrCodeExpired
	}

	return userId, nil
}

func DeleteEmailVerification(userId string) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `DELETE FROM email_verification WHERE user_id = $1;`
	res, err := shopDb.Exec(query, userId)

	if err != nil {
		return fmt.Errorf("failed to delete mail verification: <%w>", err)
	}

	if rows, err := res.RowsAffected(); rows == 0 || err != nil {
		if err != nil {
			return fmt.Errorf("failed to get rows affected %s: <%w>", userId, err)
		} else {
			return fmt.Errorf("failed to delete mail verification %s: %d rows affected, expected at least 1", userId, rows)
		}
	}

	return nil
}
