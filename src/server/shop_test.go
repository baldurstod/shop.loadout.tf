package main_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"testing"

	"shop.loadout.tf/src/server/api"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/databases/printfuldb"
	"shop.loadout.tf/src/server/databases/shop"
	"shop.loadout.tf/src/server/email"
	"shop.loadout.tf/src/server/kmip"
	"shop.loadout.tf/src/server/model"
	"shop.loadout.tf/src/server/printful"
	"shop.loadout.tf/src/server/release"
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
	shop.InitShopDB(testConfig.Databases.Shop)
	kmip.InitKmip(testConfig.Kms)
	return nil
}

/*
func RefreshAllProducts() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		printful.RefreshAllProducts("USD", true)
	}()
	wg.Wait()
}
*/

/*
func TestRefreshAllProducts(t *testing.T) {
	RefreshAllProducts()
}
*/

var username = "test@example.com"
var userPass = "test_pass"

func TestCreateUser(t *testing.T) {
	user, err := shop.CreateUser(username, userPass)
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
	email.SetMailConfig(testConfig.SMTP)
	if err := email.SendMail("noreply@loadout.tf", "noreply@loadout.tf",
		"A very very long\n  subject header spanning multiple lines",
		"test test\n\nMore test text", nil); err != nil {
		t.Error(err)
		return
	}
}

func TestSendMailHtml(t *testing.T) {
	email.SetMailConfig(testConfig.SMTP)
	if err := email.SendMailHtml("noreply@loadout.tf", "noreply@loadout.tf",
		"A very very long\n  subject header spanning multiple lines",
		`
	<html>
	<body>
	<h1>test</h1>
	</body>
	</html>
	`, nil); err != nil {
		t.Error(err)
		return
	}
}

func TestSendOrderMail(t *testing.T) {
	// Force test mode to have test links
	release.ReleaseMode = "false"

	email.SetMailConfig(testConfig.SMTP)
	user, err := shop.FindUserByID("O4LML1O5X7B2")
	if err != nil {
		t.Error(err)
		return
	}

	order, err := shop.GetOrder("9WXBA9BKM16V")
	if err != nil {
		t.Error(err)
		return
	}

	if err := email.SendOrderMail(*user, *order); err != nil {
		t.Error(err)
		return
	}
}

func TestCreateUser2(t *testing.T) {
	user := model.NewUser()
	user.AddOrder("a")
	user.Currency = "d"

	user.Cart.AddQuantity("sffds", 1)

	fmt.Printf("%+v\n", user)
}

func TestAttachOrder(t *testing.T) {
	user, err := api.GetUser(username, userPass)
	if err != nil {
		t.Error(err)
		return
	}

	//fields := shop.UpdateUserFields{}
	//fields.AddOrder = "test_order"

	err = shop.UserAddOrder(user.ID, "test_order")
	if err != nil {
		t.Error(err)
		return
	}
}
