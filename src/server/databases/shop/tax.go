package shop

import (
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

func SetTaxRate(countryCode string, stateCode string, postalCode string, city string, rate decimal.Decimal) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	_, err := shopDb.Exec(`INSERT INTO tax (country_code, state_code, postal_code, city, rate)
						VALUES ($1, $2, $3, $4, $5)`,
		countryCode,
		stateCode,
		postalCode,
		city,
		rate,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to insert tax rate: <%w>", err)
	}

	return nil
}

func GetTaxRate(countryCode string, stateCode string, postalCode string, city string) (*decimal.Decimal, error) {
	if shopDb == nil {
		return nil, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT rate FROM tax WHERE country_code = $1 AND state_code = $2 AND postal_code = $3 AND city = $4;`
	row := shopDb.QueryRow(query, countryCode, stateCode, postalCode, city)

	var rate string

	err := row.Scan(&rate)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row in GetTaxRate: <%w>", err)
	}

	r, err := decimal.NewFromString(rate)
	if err != nil {
		return nil, fmt.Errorf("failed to decode price in GetTaxRate: <%w>", err)
	}

	return &r, nil
}
