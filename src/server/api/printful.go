package api

import (
	"errors"
	"log"
	"regexp"
	"strconv"

	"github.com/baldurstod/go-printful-api-model/schemas"
	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	printfulapi "shop.loadout.tf/src/server/api/printful"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/model"
	"shop.loadout.tf/src/server/mongo/printfuldb"
	"shop.loadout.tf/src/server/printful"
	//"shop.loadout.tf/src/server/sessions"
)

var printfulConfig config.Printful

var IsAlphaNumeric = regexp.MustCompile(`^[0-9a-zA-Z]+$`).MatchString

func getCountries(c *gin.Context) error {
	countries, err := printfuldb.FindCountries()
	if err != nil {
		return err
	}

	jsonSuccess(c, map[string]interface{}{"countries": countries})

	return nil
}

func computeTaxRate(order *model.Order) error {
	recipient := schemas.TaxAddressInfo{
		City:        order.ShippingAddress.City,
		CountryCode: order.ShippingAddress.CountryCode,
		StateCode:   order.ShippingAddress.StateCode,
		ZIP:         order.ShippingAddress.PostalCode,
	}

	taxInfo, err := printful.CalculateTaxRate(recipient)
	if err != nil {
		log.Println(err)
		return errors.New("error while calling printful api")
	}

	order.TaxInfo.Required = taxInfo.Required
	order.TaxInfo.Rate = taxInfo.Rate
	order.TaxInfo.ShippingTaxable = taxInfo.ShippingTaxable

	return nil
}

func createPrintfulOrder(order *model.Order) error {
	recipient := printfulmodel.Address{
		Address1:    order.ShippingAddress.Address1,
		Address2:    order.ShippingAddress.Address2,
		City:        order.ShippingAddress.City,
		CountryCode: order.ShippingAddress.CountryCode,
		StateCode:   order.ShippingAddress.StateCode,
		ZIP:         order.ShippingAddress.PostalCode,
	}

	orderItems := make([]printfulmodel.CatalogItem, 0, len(order.Items))

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

		log.Println("AAAAAAAAAAAAAAAAAAAAAA", orderItem.RetailPrice.String())
		orderItems = append(orderItems, item)
	}

	_, err := printful.CreateOrder(order.ID, order.ShippingMethod, recipient, orderItems, nil, nil)
	if err != nil {
		log.Println(err)
		return errors.New("error while creating printful order")
	}

	return nil
}

func productToPlacementList(p *model.Product) (printfulmodel.PlacementsList, error) {
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
	}

	return placementsList, nil
}

func apiGetPrintfulProducts(c *gin.Context, params map[string]interface{}) error {
	var currency string
	if c, ok := params["currency"]; ok {
		c2, ok := c.(string)
		if ok {
			currency = c2
		}
	}

	products, err := printfulapi.GetProducts()
	if err != nil {
		return err
	}

	var variantsPrices []printfulapi.VariantPrice
	if currency != "" {
		variantsPrices, err = printfulapi.GetVariantsPrices(currency, printfulConfig.Markup)
		if err != nil {
			return err
		}
	}

	jsonSuccess(c, map[string]any{
		"products": products,
		"prices":   variantsPrices,
	})

	return nil
}

func apiGetPrintfulProduct(c *gin.Context, params map[string]interface{}) error {
	productID := int(params["product_id"].(float64))
	product, err := printfulapi.GetProduct(productID)

	if err != nil {
		return err
	}

	variants, err := printfulapi.GetVariants(productID)

	if err != nil {
		return err
	}

	jsonSuccess(c, map[string]interface{}{
		"product":  product,
		"variants": variants,
	})

	return nil
}

func apiGetPrintfulCategories(c *gin.Context) error {
	categories, err := printfulapi.GetCategories()

	if err != nil {
		return err
	}

	jsonSuccess(c, categories)

	return nil
}

func apiGetPrintfulMockupTemplates(c *gin.Context, params map[string]interface{}) error {
	templates, err := printfulapi.GetMockupTemplates(int(params["product_id"].(float64)))
	log.Println(params)

	if err != nil {
		return err
	}

	jsonSuccess(c, map[string]interface{}{
		"templates": templates,
	})

	return nil
}

func apiGetPrintfulMockupStyles(c *gin.Context, params map[string]interface{}) error {
	styles, err := printfulapi.GetMockupStyles(int(params["product_id"].(float64)))
	log.Println(params)

	if err != nil {
		return err
	}

	jsonSuccess(c, map[string]interface{}{
		"styles": styles,
	})

	return nil
}

func apiGetPrintfulProductPrices(c *gin.Context, params map[string]interface{}) error {
	productID := int(params["product_id"].(float64))
	currency := params["currency"].(string)

	prices, err := printfulapi.GetProductPrices(productID, currency, printfulConfig.Markup)

	if err != nil {
		return err
	}

	jsonSuccess(c, prices)

	return nil
}
