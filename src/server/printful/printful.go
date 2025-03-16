package printful

import (
	"log"

	printfulsdk "github.com/baldurstod/go-printful-sdk"
	"shop.loadout.tf/src/server/config"
)

var printfulConfig config.Printful
var printfulClient *printfulsdk.PrintfulClient = printfulsdk.NewPrintfulClient("")

func SetPrintfulConfig(config config.Printful) {
	printfulConfig = config
	log.Println(config)
	printfulClient.SetAccessToken(config.AccessToken)
	//go initAllProducts()
}
