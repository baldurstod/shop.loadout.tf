package api

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/baldurstod/go-printful-api-model/schemas"
	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"github.com/gin-gonic/gin"
	printfulapi "shop.loadout.tf/src/server/api/printful"
	"shop.loadout.tf/src/server/databases/printfuldb"
	"shop.loadout.tf/src/server/logger"
	"shop.loadout.tf/src/server/model"
	"shop.loadout.tf/src/server/printful"
)

var markup = 20.0

func SetMarkup(m float64) {
	if m < 20 {
		return
	}
	markup = m
}

var IsAlphaNumeric = regexp.MustCompile(`^[0-9a-zA-Z]+$`).MatchString

func apiGetCountries(c *gin.Context) apiError {
	countries, err := printfuldb.FindCountries()
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
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
		return fmt.Errorf("error while calculating tax rate: %w", err)
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
		item := printfulmodel.NewCatalogItem()

		variantID, err := strconv.Atoi(orderItem.Product.ExternalID1)
		if err != nil {
			return fmt.Errorf("error while creating printful order: %w", err)
		}
		item.CatalogVariantID = variantID

		item.ID = id
		item.ExternalID = orderItem.Product.ID
		item.Quantity = int(orderItem.Quantity)
		item.RetailPrice = orderItem.RetailPrice.String()
		item.Name = orderItem.Name
		_, item.Placements, err = orderItem.Product.GetPlacementList()
		if err != nil {
			return fmt.Errorf("error while creating printful order: %w", err)
		}

		orderItems = append(orderItems, item)
	}

	_, err := printful.CreateOrder(order.ID, order.ShippingMethod, recipient, orderItems, nil, nil)
	if err != nil {
		return fmt.Errorf("error while creating printful order: %w", err)
	}

	return nil
}

func apiGetPrintfulProducts(c *gin.Context, params map[string]any) apiError {
	var currency string
	if c, ok := params["currency"].(string); ok {
		currency = c
	}

	products, err := printfulapi.GetProducts()
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	var variantsPrices []printfulapi.VariantPrice
	if currency != "" {
		variantsPrices, err = printfulapi.GetVariantsPrices(currency, markup)
		if err != nil {
			logger.Log(c, err)
			return CreateApiError(UnexpectedError)
		}
	}

	jsonSuccess(c, map[string]any{
		"products": products,
		"prices":   variantsPrices,
	})

	return nil
}

func apiGetPrintfulProduct(c *gin.Context, params map[string]any) apiError {
	id, ok := params["product_id"].(float64)
	if !ok {
		return CreateApiError(InvalidParamProductID)
	}

	productID := int(id)

	product, err := printfulapi.GetProduct(productID)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	variants, err := printfulapi.GetVariants(productID)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, map[string]any{
		"product":  product,
		"variants": variants,
	})

	return nil
}

func apiGetPrintfulCategories(c *gin.Context) apiError {
	categories, err := printfulapi.GetCategories()
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, map[string]any{"categories": categories})

	return nil
}

func apiGetPrintfulMockupTemplates(c *gin.Context, params map[string]any) apiError {
	productID, ok := params["product_id"].(float64)
	if !ok {
		return CreateApiError(InvalidParamProductID)
	}

	templates, err := printfulapi.GetMockupTemplates(int(productID))
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, map[string]any{
		"templates": templates,
	})

	return nil
}

func apiGetPrintfulMockupStyles(c *gin.Context, params map[string]any) apiError {
	productID, ok := params["product_id"].(float64)
	if !ok {
		return CreateApiError(InvalidParamProductID)
	}

	styles, err := printfulapi.GetMockupStyles(int(productID))
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, map[string]any{
		"styles": styles,
	})

	return nil
}

func apiGetPrintfulProductPrices(c *gin.Context, params map[string]any) apiError {
	productID, ok := params["product_id"].(float64)
	if !ok {
		return CreateApiError(InvalidParamProductID)
	}

	currency, ok := params["currency"].(string)
	if !ok {
		return CreateApiError(InvalidParamCurrency)
	}

	prices, err := printfulapi.GetProductPrices(int(productID), currency, markup)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, map[string]any{"prices": prices})

	return nil
}
