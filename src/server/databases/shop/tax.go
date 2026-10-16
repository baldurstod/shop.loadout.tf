package shop

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

func SetTaxRate(countryCode string, stateCode string, postalCode string, city string, rate decimal.Decimal) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	if (countryCode == "US") && (stateCode == "" || postalCode == "" || city == "") {
		return fmt.Errorf("Can't insert tax rate for US, missing parameters |%s| |%s| |%s|", stateCode, postalCode, city)
	}

	if (countryCode == "CA") && (stateCode == "") {
		return fmt.Errorf("Can't insert tax rate for CA, missing state code |%s|", stateCode)
	}

	_, err := shopDb.Exec(`INSERT INTO tax (country_code, state_code, postal_code, city, rate, date_created, date_updated)
						VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		countryCode,
		stateCode,
		postalCode,
		strings.ToUpper(city),
		rate,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to insert tax rate |%s| |%s| |%s| |%s| <%w>", countryCode, stateCode, postalCode, city, err)
	}

	return nil
}

func GetTaxRate(countryCode string, stateCode string, postalCode string, city string) (*decimal.Decimal, error) {
	if shopDb == nil {
		return nil, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT rate FROM tax WHERE country_code = $1 AND state_code = $2 AND postal_code = $3 AND city = $4;`
	row := shopDb.QueryRow(query, countryCode, stateCode, postalCode, strings.ToUpper(city))

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
