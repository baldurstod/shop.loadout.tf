// #cgo windows CFLAGS: -I D:\Divers\libmongocrypt\include\mongocrypt

package main_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"sync"
	"testing"

	"shop.loadout.tf/src/server/api"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/databases"
	mongoshop "shop.loadout.tf/src/server/databases"
	"shop.loadout.tf/src/server/databases/printfuldb"
	"shop.loadout.tf/src/server/mail"
	"shop.loadout.tf/src/server/model"
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

var testConfig = config.Config{}

func initConfig() error {
	var err error
	var content []byte

	if content, err = os.ReadFile("config.json"); err != nil {
		return err
	}
	if err = json.Unmarshal(content, &testConfig); err != nil {
		return err
	}
	printful.SetPrintfulConfig(testConfig.Printful)
	printfuldb.InitPrintfulDB(testConfig.Databases.Printful)
	mongoshop.InitShopDB(testConfig.Databases.Shop)
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

var username = "test@example.com"
var userPass = "test_pass"

func TestCreateUser(t *testing.T) {
	hashedPassword, err := api.HashPassword(userPass)
	if err != nil {
		t.Error(err)
		return
	}

	user, err := databases.CreateUser(username, hashedPassword)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println("created user:", user)
}

func TestCheckPassword(t *testing.T) {
	user, err := api.GetUser(username, userPass)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println("returned user:", user)
}

func TestCheckWrongPassword(t *testing.T) {
	_, err := api.GetUser(username, "wrong_pass")
	if err == nil {
		t.Error("err is nil")
		return
	}
	if err.Error() != "wrong password" {
		t.Error(err)
		return
	}
}

func TestSendMail(t *testing.T) {
	mail.SetMailConfig(testConfig.SMTP)
	if err := mail.SendMail("noreply@loadout.tf", "noreply@loadout.tf",
		"A very very long\n  subject header spanning multiple lines",
		"test test\n\nMore test text"); err != nil {
		t.Error(err)
		return
	}
}
