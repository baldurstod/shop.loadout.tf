package main_test

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"runtime"
	"sync"
	"testing"

	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/mongo/printfuldb"
	"shop.loadout.tf/src/server/printful"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	_, filename, _, _ := runtime.Caller(0)
	// The ".." may change depending on you folder structure
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	err = initConfig()
	if err != nil {
		panic(err)
	}
}

func initConfig() error {
	var err error
	var content []byte
	config := config.Config{}

	if content, err = os.ReadFile("config.json"); err != nil {
		return err
	}
	if err = json.Unmarshal(content, &config); err != nil {
		return err
	}
	printful.SetPrintfulConfig(config.Printful)
	printfuldb.InitPrintfulDB(config.Databases.Printful)
	return nil
}

func RefreshAllProducts() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		printful.RefreshAllProducts("USD", true)
	}()
	wg.Wait()
}

func TestRefreshAllProducts(t *testing.T) {
	RefreshAllProducts()
}
