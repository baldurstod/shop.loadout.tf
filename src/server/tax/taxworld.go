package tax

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/shopspring/decimal"
	"shop.loadout.tf/src/server/databases"
	"shop.loadout.tf/src/server/databases/printfuldb"
)

func LoadWorldSalesTax(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error while opening world tax file %w", err)
	}
	defer file.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(file)
	csvReader.Comma = '\t'
	data, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("error while reading world tax file %w", err)
	}

	createWorldSalesTaxList(data)

	return nil
}

const worldNameField = 0
const worldRateField = 1

func createWorldSalesTaxList(data [][]string) error {
	countries, err := printfuldb.FindCountries()
	if err != nil {
		return err
	}

	for i, line := range data {
		if i == 0 { // omit header line
			continue
		}
		name := line[worldNameField]

		for _, j := range countries {
			if j.Name == name {
				rate, err := decimal.NewFromString(line[worldRateField])
				if err != nil {
					return fmt.Errorf("error while parsing rate %s %w", line[worldRateField], err)
				}
				if _, err := databases.SetTaxRate(j.Code, "", "", "", rate); err != nil {
					return fmt.Errorf("error while inserting tax rate %w", err)
				}
			}
		}

		/*
		 */
	}

	return nil
}
