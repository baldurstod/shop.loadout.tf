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

	jsonSuccess(c, map[string]any{"countries": countries})

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
		return errors.New("error while calculating tax rate")
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
			return errors.New("error while creating third party order")
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
			return errors.New("error while creating third party order")
		}

		log.Println("AAAAAAAAAAAAAAAAAAAAAA", orderItem.RetailPrice.String())
		orderItems = append(orderItems, item)
	}

	_, err := printful.CreateOrder(order.ID, order.ShippingMethod, recipient, orderItems, nil, nil)
	if err != nil {
		log.Println(err)
		return errors.New("error while creating third party order")
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

func apiGetPrintfulProducts(c *gin.Context, params map[string]any) error {
	var currency string
	if c, ok := params["currency"].(string); ok {
		currency = c
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

func apiGetPrintfulProduct(c *gin.Context, params map[string]any) error {
	id, ok := params["product_id"].(float64)
	if !ok {
		return errors.New("invalid product id")
	}

	productID := int(id)

	product, err := printfulapi.GetProduct(productID)
	if err != nil {
		return err
	}

	variants, err := printfulapi.GetVariants(productID)

	if err != nil {
		return err
	}

	jsonSuccess(c, map[string]any{
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

	jsonSuccess(c, map[string]any{"categories": categories})

	return nil
}

func apiGetPrintfulMockupTemplates(c *gin.Context, params map[string]any) error {
	productID, ok := params["product_id"].(float64)
	if !ok {
		return errors.New("invalid product id")
	}

	templates, err := printfulapi.GetMockupTemplates(int(productID))
	log.Println(params)

	if err != nil {
		return err
	}

	jsonSuccess(c, map[string]any{
		"templates": templates,
	})

	return nil
}

func apiGetPrintfulMockupStyles(c *gin.Context, params map[string]any) error {
	productID, ok := params["product_id"].(float64)
	if !ok {
		return errors.New("invalid product id")
	}

	styles, err := printfulapi.GetMockupStyles(int(productID))
	log.Println(params)

	if err != nil {
		return err
	}

	jsonSuccess(c, map[string]any{
		"styles": styles,
	})

	return nil
}

func apiGetPrintfulProductPrices(c *gin.Context, params map[string]any) error {
	productID, ok := params["product_id"].(float64)
	if !ok {
		return errors.New("invalid product id")
	}

	currency, ok := params["currency"].(string)
	if !ok {
		return errors.New("invalid currency")
	}

	prices, err := printfulapi.GetProductPrices(int(productID), currency, printfulConfig.Markup)

	if err != nil {
		return err
	}

	jsonSuccess(c, map[string]any{"prices": prices})

	return nil
}
