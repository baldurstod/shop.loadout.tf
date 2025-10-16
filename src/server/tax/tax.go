package tax

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shopspring/decimal"
	"shop.loadout.tf/src/server/databases"
)

func LoadUSTax(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error while opening US tax file %w", err)
	}
	defer file.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(file)
	csvReader.Comma = '\t'
	data, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("error while reading US tax file %w", err)
	}

	createUSTaxList(data)

	return nil
}

const postalCodeField = 0
const stateField = 1
const countyField = 2
const cityField = 3
const taxField = 4

func createUSTaxList(data [][]string) error {
	exist := map[string]map[string]map[string]bool{}
	skippedUnincorporated := 0
	skippedDuplicate := 0

	for i, line := range data {
		if i == 0 { // omit header line
			continue
		}

		rate, err := decimal.NewFromString(line[taxField]) //strconv.ParseFloat(line[taxField], 32)
		if err != nil {
			return fmt.Errorf("error while parsing rate %s %w", line[taxField], err)
		}

		state := line[stateField]
		postalCode := line[postalCodeField]
		city := line[cityField]

		if strings.HasPrefix(strings.ToUpper(city), "UNINCORPORATED") {
			skippedUnincorporated += 1
			continue
		}

		s1, found := exist[state]
		if !found {
			s1 = map[string]map[string]bool{}
			exist[state] = s1
		}

		s2, found := s1[postalCode]
		if !found {
			s2 = map[string]bool{}
			s1[postalCode] = s2
		}

		_, found = s2[city]
		if found {
			skippedDuplicate += 1
			continue
		} else {
			s2[city] = true
		}

		if _, err := databases.SetTaxRate("US", state, postalCode, city, rate); err != nil {
			return fmt.Errorf("error while inserting tax rate %w", err)
		}
	}

	log.Println("Unincorporated:", skippedUnincorporated, "Duplicate:", skippedDuplicate)

	return nil
}
