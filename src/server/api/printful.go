package api

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	printfulApiModel "github.com/baldurstod/go-printful-api-model"
	"github.com/baldurstod/go-printful-api-model/requests"
	"github.com/baldurstod/go-printful-api-model/responses"
	"github.com/baldurstod/go-printful-api-model/schemas"
	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/model"
	"shop.loadout.tf/src/server/mongo"

	//"shop.loadout.tf/src/server/sessions"
	"bytes"

	"github.com/gin-contrib/sessions"
)

var printfulConfig config.Printful
var printfulURL string
var printfulClient = createPrinfulClient()

var IsAlphaNumeric = regexp.MustCompile(`^[0-9a-zA-Z]+$`).MatchString

func createPrinfulClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{},
	}
}

func SetPrintfulConfig(config config.Printful) {
	printfulConfig = config
	log.Println(config)
	var err error
	printfulURL, err = url.JoinPath(printfulConfig.Endpoint, "/api")
	if err != nil {
		panic("Error while getting printful url")
	}

	u, err := url.Parse(printfulConfig.Endpoint)
	if err == nil {
		if strings.HasPrefix(u.Host, "localhost") {
			printfulClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
	}
}

func fetchAPI(action string, version int, params interface{}) (*http.Response, error) {

	body := map[string]interface{}{
		"action":  action,
		"version": version,
		"params":  params,
	}

	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Fetching printful api %s version %d \n", action, version)
	res, err := printfulClient.Post(printfulURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	fmt.Println("Printful api returned code", res.StatusCode)

	return res, err
}

func getCountries(c *gin.Context) error {
	/*
		u, err := url.JoinPath(printfulConfig.Endpoint, "/countries")
		if err != nil {
			return errors.New("error while getting printful url")
		}

		resp, err := http.Get(u)*/
	resp, err := fetchAPI("get-countries", 1, nil)
	//body, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(body))
	log.Println(resp)

	if err != nil {
		log.Println(err)
		return errors.New("error while calling printful api")
	}
	defer resp.Body.Close()

	countriesResponse := printfulApiModel.CountriesResponse{}
	err = json.NewDecoder(resp.Body).Decode(&countriesResponse)
	if err != nil {
		log.Println(err)
		return errors.New("error while decoding printful response")
	}

	jsonSuccess(c, map[string]interface{}{"countries": countriesResponse.Countries})

	return nil
}

func computeTaxRate(order *model.Order) error {
	calculateTaxRates := requests.CalculateTaxRate{
		Recipient: schemas.TaxAddressInfo{
			City:        order.ShippingAddress.City,
			CountryCode: order.ShippingAddress.CountryCode,
			StateCode:   order.ShippingAddress.StateCode,
			ZIP:         order.ShippingAddress.PostalCode,
		},
	}
	resp, err := fetchAPI("calculate-tax-rate", 1, calculateTaxRates)
	if err != nil {
		log.Println(err)
		return errors.New("error while calling printful api")
	}
	defer resp.Body.Close()

	response := responses.TaxRates{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println(err)
		return errors.New("error while decoding printful response")
	}

	order.TaxInfo.Required = response.Result.Required
	order.TaxInfo.Rate = response.Result.Rate
	order.TaxInfo.ShippingTaxable = response.Result.ShippingTaxable

	log.Println(response)
	return nil
}

type calculateShippingRatesResponse struct {
	Success       bool                         `json:"success"`
	ShippingInfos []printfulmodel.ShippingRate `json:"result"`
}

func apiCalculateShippingRates(c *gin.Context, s sessions.Session) error {
	orderID, ok := s.Get("order_id").(string)
	if !ok {
		return errors.New("error while retrieving order id")
	}
	order, err := mongo.FindOrder(orderID)
	if err != nil {
		log.Println(err)
		return errors.New("error while retrieving order")
	}

	calculateShippingRatesRequest := printfulApiModel.CalculateShippingRatesRequest{Items: []printfulmodel.CatalogOrWarehouseShippingRateItem{}}
	calculateShippingRatesRequest.Recipient.Address1 = order.ShippingAddress.Address1
	calculateShippingRatesRequest.Recipient.City = order.ShippingAddress.City
	calculateShippingRatesRequest.Recipient.CountryCode = order.ShippingAddress.CountryCode
	calculateShippingRatesRequest.Recipient.StateCode = order.ShippingAddress.StateCode
	calculateShippingRatesRequest.Recipient.ZIP = order.ShippingAddress.PostalCode

	for _, orderItem := range order.Items {
		p, err := mongo.GetProduct(orderItem.ProductID)
		if err != nil {
			log.Println(err)
			return errors.New("error while computing shipping rates")
		}

		variantID, err := strconv.Atoi(p.ExternalID1)
		if err != nil {
			log.Println(err)
			return errors.New("error while computing shipping rates")
		}

		itemInfo := printfulmodel.CatalogOrWarehouseShippingRateItem{
			Source:           "catalog",
			CatalogVariantID: variantID,
			Quantity:         int(orderItem.Quantity),
		}

		calculateShippingRatesRequest.Items = append(calculateShippingRatesRequest.Items, itemInfo)
	}

	/*order.ShippingAddress = address
	err = mongo.UpdateOrder(order)
	if err != nil {
		log.Println(err)
		return errors.New("error while updating order")
	}

	log.Println(order)*/
	jsonSuccess(c, map[string]interface{}{"order": order})
	return nil
}

/*
type CalculateShippingRatesRequest struct {
	Recipient AddressInfo `json:"recipient" bson:"recipient" mapstructure:"recipient"`
	Items     []ItemInfo  `json:"items" bson:"items" mapstructure:"items"`
	Currency  string      `json:"currency" bson:"currency" mapstructure:"currency"`
	Locale    string      `json:"locale" bson:"locale" mapstructure:"locale"`
}

type AddressInfo struct {
	Address1    string `json:"address1" bson:"address1" mapstructure:"address1"`
	City        string `json:"city" bson:"city" mapstructure:"city"`
	CountryCode string `json:"country_code" bson:"country_code" mapstructure:"country_code"`
	StateCode   string `json:"state_code" bson:"state_code" mapstructure:"state_code"`
	ZIP         string `json:"zip" bson:"zip" mapstructure:"zip"`
	Phone       string `json:"phone" bson:"phone" mapstructure:"phone"`
}
type ItemInfo struct {
	VariantID                 string `json:"variant_id" bson:"variant_id" mapstructure:"variant_id"`
	ExternalVariantID         string `json:"external_variant_id" bson:"external_variant_id" mapstructure:"external_variant_id"`
	WarehouseProductVariantID string `json:"warehouse_product_variant_id" bson:"warehouse_product_variant_id" mapstructure:"warehouse_product_variant_id"`
	Quantity                  int    `json:"quantity" bson:"quantity" mapstructure:"quantity"`
	Value                     string `json:"value" bson:"value" mapstructure:"value"`
}


type Order struct {
	ID                 primitive.ObjectID      `json:"id" bson:"_id"`
	Currency           string                  `json:"currency" bson:"currency"`
	DateCreated        int64                   `json:"date_created" bson:"date_created"`
	DateUpdated        int64                   `json:"date_updated" bson:"date_updated"`
	ShippingAddress    Address                 `json:"shipping_address" bson:"shipping_address"`
	BillingAddress     Address                 `json:"billing_address" bson:"billing_address"`
	SameBillingAddress bool                    `json:"same_billing_address" bson:"same_billing_address"`
	Items              []OrderItem             `json:"items" bson:"items"`
	ShippingInfos      map[string]ShippingInfo `json:"shipping_infos" bson:"shipping_infos"`
	TaxInfo            TaxInfo                 `json:"tax_info" bson:"tax_info"`
	ShippingMethod     string                  `json:"shipping_method" bson:"shipping_method"`
	PrintfulOrderID    string                  `json:"printful_order_id" bson:"printful_order_id"`
	PaypalOrderID      string                  `json:"paypal_order_id" bson:"paypal_order_id"`
	Status             string                  `json:"status" bson:"status"`
}
type OrderItem struct {
	ProductID    string  `json:"product_id" bson:"product_id"`
	Name         string  `json:"name" bson:"name"`
	Quantity     uint    `json:"quantity" bson:"quantity"`
	RetailPrice  float64 `json:"retail_price" bson:"retail_price"`
	ThumbnailURL string  `json:"thumbnail_url" bson:"thumbnail_url"`
}
type Address struct {
	FirstName    string `json:"first_name" bson:"first_name" mapstructure:"first_name"`
	LastName     string `json:"last_name" bson:"last_name" mapstructure:"last_name"`
	Organization string `json:"organization" bson:"organization" mapstructure:"organization"`
	Address1     string `json:"address1" bson:"address1" mapstructure:"address1"`
	Address2     string `json:"address2" bson:"address2" mapstructure:"address2"`
	City         string `json:"city" bson:"city" mapstructure:"city"`
	StateCode    string `json:"state_code" bson:"state_code" mapstructure:"state_code"`
	StateName    string `json:"state_name" bson:"state_name" mapstructure:"state_name"`
	CountryCode  string `json:"country_code" bson:"country_code" mapstructure:"country_code"`
	CountryName  string `json:"country_name" bson:"country_name" mapstructure:"country_name"`
	PostalCode   string `json:"postal_code" bson:"postal_code" mapstructure:"postal_code"`
	Phone        string `json:"phone" bson:"phone" mapstructure:"phone"`
	Email        string `json:"email" bson:"email" mapstructure:"email"`
	TaxNumber    string `json:"tax_number" bson:"tax_number" mapstructure:"tax_number"`
}

*/

type createPrintfulOrderResponse struct {
	Success bool          `json:"success"`
	Order   schemas.Order `json:"result"`
}

func createPrintfulOrder(order *model.Order) error {
	printfulOrder := requests.CreateOrder{
		ExternalID: order.ID,
		Shipping:   order.ShippingMethod,
		Recipient: printfulmodel.Address{
			Address1:    order.ShippingAddress.Address1,
			Address2:    order.ShippingAddress.Address2,
			City:        order.ShippingAddress.City,
			CountryCode: order.ShippingAddress.CountryCode,
			StateCode:   order.ShippingAddress.StateCode,
			ZIP:         order.ShippingAddress.PostalCode,
		},
		OrderItems: make([]printfulmodel.CatalogItem, 0, len(order.Items)),
	}

	log.Println(printfulOrder)

	for id, orderItem := range order.Items {
		log.Println("**********************", orderItem)
		item := printfulmodel.NewCatalogItem()

		variantID, err := strconv.Atoi(orderItem.Product.ExternalID1)
		if err != nil {
			log.Println(err)
			return errors.New("error while creating printful order")
		}
		item.CatalogVariantID = variantID

		item.ID = id
		item.ExternalID = orderItem.Product.ID
		item.Quantity = int(orderItem.Quantity)
		item.RetailPrice = orderItem.RetailPrice.String()
		item.Name = orderItem.Name
		item.Placements, err = productToPlacementList(&orderItem.Product)
		if err != nil {
			log.Println(err)
			return errors.New("error while creating printful order")
		}

		/*
			{
				ExternalVariantID: orderItem.ProductID,
				Quantity:          int(orderItem.Quantity),
				RetailPrice:       orderItem.RetailPrice.String(),
			}
		*/
		log.Println("AAAAAAAAAAAAAAAAAAAAAA", orderItem.RetailPrice.String())
		printfulOrder.OrderItems = append(printfulOrder.OrderItems, item)
	}

	resp, err := fetchAPI("create-order", 1, printfulOrder)
	if err != nil {
		log.Println(err)
		return errors.New("error while calling printful api")
	}
	defer resp.Body.Close()

	response := responses.CreateOrderResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println(err)
		return errors.New("error while decoding printful response")
	}

	if !response.Success {
		log.Println(response)
		return errors.New("error while creating printful order")
	}

	//jsonSuccess(c, map[string]interface{}{"order": response.Order})

	return nil
}

func productToPlacementList(p *model.Product) (printfulmodel.PlacementsList, error) {
	//extraData := map[string]any{"printful": map[string]any{"placements": extraDataPlacements}}

	productExtraData := model.ProductExtraData{}
	err := mapstructure.Decode(p.ExtraData, &productExtraData)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error while reading params")
	}

	log.Println(productExtraData)

	placementsList := make(printfulmodel.PlacementsList, len(productExtraData.Printful.Placements))

	for i, placement := range productExtraData.Printful.Placements {
		placementsList[i] = printfulmodel.Placement{
			Placement:     placement.Placement,
			Technique:     placement.Technique,
			PrintAreaType: "simple", //TODO: variable ?
			Layers: []printfulmodel.Layer{{
				Type: "file", //TODO: variable ?
				Url:  placement.ImageURL,
				//LayerOptions
				//LayerPosition
			}},
		}
		/*

			Placement         string  `json:"placement" bson:"placement"`
			Technique         string  `json:"technique" bson:"technique"`
			PrintAreaType     string  `json:"print_area_type" bson:"print_area_type"`
			Layers            []Layer `json:"layers" bson:"layers"`
			PlacementOptions  `json:"placement_options" bson:"placement_options"`
			Status            string `json:"status" bson:"status"`
			StatusExplanation string `json:"status_explanation" bson:"status_explanation"`
		*/
	}

	/*

		Placement   string `json:"placement" bson:"placement" mapstructure:"placement"`
		Technique   string `json:"technique" bson:"technique" mapstructure:"technique"`
		Orientation string `json:"orientation" bson:"orientation" mapstructure:"orientation"`
		ImageURL    string `json:"image_url" bson:"image_url" mapstructure:"image_url"`
		ThumbURL    string `json:"thumb_url" bson:"thumb_url" mapstructure:"thumb_url"`
	*/

	return placementsList, nil
}

/*
roundPrice(currency, price) {
	let digits = CURRENCIES_DIGITS[currency] ?? 2;
	return Number(Number.parseFloat(price).toFixed(digits));
}
*/
