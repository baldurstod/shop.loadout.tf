package printfuldb

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

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

func GetCountryName(countryCode string, stateCode string) (string, string, error) {
	if printfulDb == nil {
		return "", "", errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	var countryName string
	var stateName string
	var states string

	query := `SELECT name, states FROM countries WHERE code = $1;`
	row := printfulDb.QueryRow(query, countryCode)

	err := row.Scan(&countryName, &states)
	if err != nil {
		return "", "", fmt.Errorf("failed to scan row in GetCountryName: <%w>", err)
	}

	statesJson := []printfulmodel.State{}
	if err = json.Unmarshal([]byte(states), &statesJson); err != nil {
		return "", "", err
	}
	if stateCode != "" {
		idx := slices.IndexFunc(statesJson, func(c printfulmodel.State) bool { return c.Code == stateCode })
		if idx == -1 {
			return "", "", fmt.Errorf("failed to get state name %s GetCountryName", stateCode)
		}
		stateName = statesJson[idx].Name
	}

	return countryName, stateName, nil
}
