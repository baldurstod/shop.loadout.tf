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
			api.SetPrintfulConfig(config.Printful)
			mongo.InitMongoDB(config.Database)
			server.StartServer(config.HTTP)
		} else {
			log.Println("Error while reading configuration", err)
		}
	} else {
		log.Println("Error while reading configuration file", err)
	}
}
