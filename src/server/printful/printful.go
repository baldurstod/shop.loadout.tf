package printful

import (
	printfulsdk "github.com/baldurstod/go-printful-sdk"
	"shop.loadout.tf/src/server/config"
)

var printfulConfig config.Printful
var printfulClient *printfulsdk.PrintfulClient = printfulsdk.NewPrintfulClient("")

func SetPrintfulConfig(config config.Printful) {
	printfulConfig = config
	printfulClient.SetAccessToken(config.AccessToken)
	//go initAllProducts()
}

func GetMarkup() float64 {
	return printfulConfig.Markup
}
