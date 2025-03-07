package main

import (
	"encoding/json"
	"log"
	"os"

	"shop.loadout.tf/src/server/api"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/mongo"
	"shop.loadout.tf/src/server/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config := config.Config{}
	if content, err := os.ReadFile("config.json"); err == nil {
		if err = json.Unmarshal(content, &config); err == nil {
			api.SetImagesConfig(config.Images)
			api.SetPrintfulConfig(config.Printful)
			api.SetPaypalConfig(config.Paypal)
			mongo.InitShopDB(config.Databases.Shop)
			mongo.InitImagesDB(config.Databases.Images)
			go api.RunTasks()
			server.StartServer(config)
		} else {
			log.Println("Error while reading configuration", err)
		}
	} else {
		log.Println("Error while reading configuration file", err)
	}
}
