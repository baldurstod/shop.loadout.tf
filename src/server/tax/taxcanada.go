package tax

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/shopspring/decimal"
	"shop.loadout.tf/src/server/databases"
)

func LoadCanadaTax(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error while opening Canada tax file %w", err)
	}
	defer file.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(file)
	csvReader.Comma = '\t'
	data, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("error while reading Canada tax file %w", err)
	}

	createCanadaTaxList(data)

	return nil
}

const canadaState = 1
const canadaRate = 2

func createCanadaTaxList(data [][]string) error {
	for i, line := range data {
		if i == 0 { // omit header line
			continue
		}

		rate, err := decimal.NewFromString(line[canadaRate])
		if err != nil {
			return fmt.Errorf("error while parsing rate %s %w", line[taxField], err)
		}

		state := line[canadaState]

		if _, err := databases.SetTaxRate("CA", state, "", "", rate); err != nil {
			return fmt.Errorf("error while inserting tax rate %w", err)
		}
	}

	return nil
}
