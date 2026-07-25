package shop

import (
	"errors"
	"fmt"
	"time"
)

func insertKek(key []byte, nonce []byte, encryptionTag []byte) (kekId int64, err error) {
	if shopDb == nil {
		return 0, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	err = shopDb.QueryRow(`INSERT INTO keks (key, nonce, authenticated_encryption_tag, date_created) VALUES ($1, $2, $3, $4) RETURNING id`,
		key,
		nonce,
		encryptionTag,
		time.Now(),
	).Scan(&kekId)

	if err != nil {
		return 0, fmt.Errorf("failed to insert key: <%w>", err)
	}

	return kekId, nil
}

func getKek(kekId int64) ([]byte, []byte, []byte, error) {
	if shopDb == nil {
		return nil, nil, nil, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT key, nonce, authenticated_encryption_tag FROM keks WHERE id = $1;`

	row := shopDb.QueryRow(query, kekId)

	var key []byte
	var nonce []byte
	var encryptionTag []byte

	err := row.Scan(&key, &nonce, &encryptionTag)
	if err != nil {
		return nil, nil, nil, err
	}

	return key, nonce, encryptionTag, nil
}
