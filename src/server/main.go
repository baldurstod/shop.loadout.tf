package main

import (
	"encoding/json"
	"log"
	"os"

	"shop.loadout.tf/src/server/api"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/databases/postgre"
	"shop.loadout.tf/src/server/databases/printfuldb"
	"shop.loadout.tf/src/server/databases/shop"
	"shop.loadout.tf/src/server/kmip"
	"shop.loadout.tf/src/server/printful"
	"shop.loadout.tf/src/server/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config := config.Config{}

	if content, err := os.ReadFile("config.json"); err == nil {
		if err = json.Unmarshal(content, &config); err == nil {
			api.SetImagesConfig(config.Images)
			api.SetPaypalConfig(config.Paypal)
			printful.SetPrintfulConfig(config.Printful)
			shop.InitShopDB(config.Databases.Shop)
			kmip.InitKmip(config.Kms)
			printfuldb.InitPrintfulDB(config.Databases.Printful)
			server.InitsessionsDB(config.Sessions.DB)
			api.SetMarkup(printful.GetMarkup())
			api.RunTasks()
			server.StartServer(config)
			defer postgre.ClosePostgre()
		} else {
			log.Println("Error while reading configuration", err)
		}
	} else {
		log.Println("Error while reading configuration file", err)
	}
}
