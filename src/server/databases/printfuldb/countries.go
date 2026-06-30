package printfuldb

import (
	"encoding/json"
	"errors"
	"fmt"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
)

func FindCountries() ([]printfulmodel.Country, error) {
	if printfulDb == nil {
		return nil, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT code, name, region, states FROM countries;`
	res, err := printfulDb.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query "+query+"in FindCountries: <%w>", err)
	}
	defer res.Close()

	countries := make([]printfulmodel.Country, 0, 200)
	for res.Next() {
		var name string
		var code string
		var region string
		var states string

		err = res.Scan(&code, &name, &region, &states)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row in FindCountries: <%w>", err)
		}

		statesJson := []printfulmodel.State{}
		if err = json.Unmarshal([]byte(states), &statesJson); err != nil {
			return nil, err
		}

		country := printfulmodel.Country{Name: name, Code: code, Region: region, States: statesJson}

		countries = append(countries, country)
	}

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("failed to get next row in FindCountries: <%w>", err)
	}

	return countries, nil
}
