package main

import (
	"encoding/json"
	"log"
	"os"

	"shop.loadout.tf/src/server/api"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/mongo"
	"shop.loadout.tf/src/server/mongo/printfuldb"
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
			mongo.InitShopDB(config.Databases.Shop)
			mongo.InitImagesDB(config.Databases.Images)
			printfuldb.InitPrintfulDB(config.Databases.Printful)
			go api.RunTasks()
			server.StartServer(config)
			defer mongo.Cleanup()
		} else {
			log.Println("Error while reading configuration", err)
		}
	} else {
		log.Println("Error while reading configuration file", err)
	}
}
