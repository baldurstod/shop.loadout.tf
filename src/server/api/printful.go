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

	printfulModel "github.com/baldurstod/printful-api-model"
	"github.com/baldurstod/printful-api-model/requestbodies"
	"github.com/baldurstod/printful-api-model/responses"
	"github.com/baldurstod/printful-api-model/schemas"
	"github.com/gin-gonic/gin"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/model"
	"shop.loadout.tf/src/server/model/requests"
	"shop.loadout.tf/src/server/mongo"

	//"shop.loadout.tf/src/server/sessions"
	"bytes"
	"context"
	_ "io/ioutil"
	_ "os"
	_ "strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/greatcloak/decimal"
	"github.com/mitchellh/mapstructure"
	"github.com/plutov/paypal/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var printfulConfig config.Printful
var paypalConfig config.Paypal
var printfulURL string

var IsAlphaNumeric = regexp.MustCompile(`^[0-9a-zA-Z]+$`).MatchString

func SetPrintfulConfig(config config.Printful) {
	printfulConfig = config
	log.Println(config)
	var err error
	printfulURL, err = url.JoinPath(printfulConfig.Endpoint, "/api")
	if err != nil {
		panic("Error while getting printful url")
	}
}

func SetPaypalConfig(config config.Paypal) {
	paypalConfig = config
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
	res, err := http.Post(printfulURL, "application/json", bytes.NewBuffer(requestBody))
	fmt.Println("Printful api returned code", res.StatusCode)

	return res, err
}

func getCountries(c *gin.Context) error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

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

	countriesResponse := printfulModel.CountriesResponse{}
	err = json.NewDecoder(resp.Body).Decode(&countriesResponse)
	if err != nil {
		log.Println(err)
		return errors.New("error while decoding printful response")
	}

	jsonSuccess(c, map[string]interface{}{"countries": countriesResponse.Countries})

	return nil
}

func getCurrency(c *gin.Context, s sessions.Session) error {
	jsonSuccess(c, s.Get("currency"))
	return nil
}

func getFavorites(c *gin.Context, s sessions.Session) error {
	favorites := s.Get("favorites").(map[string]interface{})

	v := make([]string, 0, len(favorites))

	for key := range favorites {
		v = append(v, key)
	}

	jsonSuccess(c, v)
	return nil
}

func getProduct(c *gin.Context, params map[string]interface{}) error {
	if params == nil {
		return errors.New("no params provided")
	}

	product, err := mongo.FindProduct(params["product_id"].(string))

	if err != nil {
		log.Println(err)
		return errors.New("error while getting products")
	}

	for _, variantID := range product.VariantIDs {
		//variants[variantID] = struct{}{}
		p, err := mongo.FindProduct(variantID)

		if err == nil {
			product.AddVariant(model.NewVariant(p))
		}
	}

	jsonSuccess(c, map[string]interface{}{"product": product})
	return nil
}

func getProducts(c *gin.Context) error {
	p, err := mongo.GetProducts()

	if err != nil {
		log.Println(err)
		return errors.New("error while getting products")
	}

	jsonSuccess(c, p)
	return nil
}

func sendContact(c *gin.Context, params map[string]interface{}) error {
	if params == nil {
		return errors.New("no params provided")
	}

	id, err := mongo.SendContact(params)

	if err != nil {
		log.Println(err)
		return errors.New("error while sending contact")
	}

	jsonSuccess(c, id)
	return nil
}

func setFavorite(c *gin.Context, s sessions.Session, params map[string]interface{}) error {
	if params == nil {
		return errors.New("no params provided")
	}

	pID, ok := params["product_id"]
	isFavorite, ok2 := params["is_favorite"]

	if !ok || !ok2 {
		return errors.New("missing params")
	}

	favorites := s.Get("favorites").(map[string]interface{})

	productId := pID.(string)
	if isFavorite.(bool) {
		favorites[productId] = struct{}{}
	} else {
		delete(favorites, productId)
	}

	log.Println(favorites)

	jsonSuccess(c, nil)
	return nil
}

func addProduct(c *gin.Context, s sessions.Session, params map[string]interface{}) error {
	if params == nil {
		return errors.New("no params provided")
	}

	pID, ok := params["product_id"]
	quantity, ok2 := params["quantity"]

	if !ok || !ok2 {
		return errors.New("missing params")
	}

	cart := s.Get("cart").(model.Cart)

	cart.AddQuantity(pID.(string), uint(quantity.(float64)))
	s.Delete("order_id")

	jsonSuccess(c, map[string]interface{}{"cart": cart})
	return nil
}

func setProductQuantity(c *gin.Context, s sessions.Session, params map[string]interface{}) error {
	if params == nil {
		return errors.New("no params provided")
	}

	pID, ok := params["product_id"]
	quantity, ok2 := params["quantity"]

	if !ok || !ok2 {
		return errors.New("missing params")
	}

	cart := s.Get("cart").(model.Cart)

	cart.SetQuantity(pID.(string), uint(quantity.(float64)))
	s.Delete("order_id")

	jsonSuccess(c, map[string]interface{}{"cart": cart})
	return nil
}

func getCart(c *gin.Context, s sessions.Session) error {
	cart := s.Get("cart").(model.Cart)

	jsonSuccess(c, map[string]interface{}{"cart": cart})
	return nil
}

func initCheckout(c *gin.Context, s sessions.Session) error {
	cart := s.Get("cart").(model.Cart)

	order, err := mongo.CreateOrder()
	if err != nil {
		log.Println(err)
		return errors.New("error while creating order")
	}
	/*
		order.ShippingAddress.FirstName = "ShippingAddress.FirstName"
		order.ShippingAddress.LastName = "ShippingAddress.LastName"
		order.ShippingAddress.Company = "ShippingAddress.Company"
		order.ShippingAddress.Address1 = "ShippingAddress.Address1"
		order.ShippingAddress.Address2 = "ShippingAddress.Address2"
		order.ShippingAddress.City = "ShippingAddress.City"
		order.ShippingAddress.StateCode = "CA"
		order.ShippingAddress.CountryCode = "US"
		order.ShippingAddress.PostalCode = "ShippingAddress.PostalCode"
		order.ShippingAddress.Phone = "ShippingAddress.Phone"
		order.ShippingAddress.Email = "ShippingAddress.Email"

		order.BillingAddress.FirstName = "ShippingAddress.FirstName"
		order.BillingAddress.LastName = "ShippingAddress.LastName"
		order.BillingAddress.Company = "ShippingAddress.Company"
		order.BillingAddress.Address1 = "ShippingAddress.Address1"
		order.BillingAddress.Address2 = "ShippingAddress.Address2"
		order.BillingAddress.City = "ShippingAddress.City"
		order.BillingAddress.StateCode = "CA"
		order.BillingAddress.CountryCode = "US"
		order.BillingAddress.PostalCode = "ShippingAddress.PostalCode"
		order.BillingAddress.Phone = "ShippingAddress.Phone"
		order.BillingAddress.Email = "ShippingAddress.Email"*/

	err = initCheckoutItems(&cart, order)
	if err != nil {
		log.Println(err)
		return errors.New("error while adding items to order")
	}

	order.Currency = cart.Currency
	now := time.Now().Unix()
	order.DateCreated = now
	order.DateUpdated = now
	order.Status = "created"

	err = mongo.UpdateOrder(order)
	if err != nil {
		log.Println(err)
		return errors.New("error while updating order")
	}

	log.Println(order)
	s.Set("order_id", order.ID.Hex())
	log.Println(s)

	jsonSuccess(c, map[string]interface{}{"order": order})

	return nil
}

func initCheckoutItems(cart *model.Cart, order *model.Order) error {
	log.Println(cart.Items)
	for productID, quantity := range cart.Items {
		p, err := mongo.GetProduct(productID)
		if err != nil {
			log.Println(err)
			return errors.New("error during order initialization")
		}

		orderItem := model.OrderItem{}
		orderItem.ProductID = p.ID.Hex()
		orderItem.Name = p.Name
		orderItem.ThumbnailURL = p.ThumbnailURL
		orderItem.Quantity = quantity
		orderItem.RetailPrice = p.RetailPrice

		order.Items = append(order.Items, orderItem)
	}

	log.Println("-----------------", order)

	return nil
}

func apiCreateProduct(c *gin.Context, s sessions.Session, params map[string]interface{}) error {
	if params == nil {
		return errors.New("no params provided")
	}
	//log.Println(params)
	//createProduct := params["product"].(requests.CreateProductRequest)

	createProductRequest := requests.CreateProductRequest{}
	err := mapstructure.Decode(params["product"], &createProductRequest)
	if err != nil {
		log.Println(err)
		return errors.New("error while reading params")
	}

	err = createProductRequest.CheckParams()
	if err != nil {
		log.Println(err)
		return errors.New("invalid params")
	}

	log.Println(createProductRequest.Name, createProductRequest.Type, createProductRequest.VariantID)
	products, err := createProduct(&createProductRequest)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error while creating product: %w", err)
	}

	jsonSuccess(c, map[string]interface{}{"products": products})

	return nil
}

func createProduct(request *requests.CreateProductRequest) ([]*model.Product, error) {
	pfVariant, err := getPrintfulVariant(request.VariantID)
	if err != nil {
		log.Println(err)
		return nil, errors.New("variant not found")
	}

	log.Println(pfVariant)
	pfProduct, err := getPrintfulProduct(pfVariant.ProductID)
	if err != nil {
		log.Println(err)
		return nil, errors.New("product not found")
	}

	log.Println(pfProduct)

	resp, err := fetchAPI("get-similar-variants", 1, map[string]interface{}{
		"variant_id": pfVariant.ID,
	})

	if err != nil {
		log.Println(err)
		return nil, errors.New("error while calling printful api")
	}

	similarVariantsResponse := printfulModel.SimilarVariantsResponse{}
	err = json.NewDecoder(resp.Body).Decode(&similarVariantsResponse)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error while decoding printful response")
	}

	if !similarVariantsResponse.Success {
		log.Println(similarVariantsResponse)
		return nil, errors.New("error while getting printful variant")
	}

	log.Println(similarVariantsResponse)

	variantCount := len(similarVariantsResponse.SimilarVariants)
	ids, err := createShopProducts(variantCount)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error while creating product: %w", err)
	}

	variants := make([]interface{}, 0, variantCount) //map[string]interface{}{}
	i := 0
	for i < variantCount {
		variant := map[string]interface{}{
			"variant_id":          similarVariantsResponse.SimilarVariants[i],
			"external_variant_id": ids[i],
			"retail_price":        9999,
		}

		variants = append(variants, variant)
		i += 1
	}

	log.Println(ids, err)

	resp, err = fetchAPI("create-sync-product", 1, map[string]interface{}{
		"product_id": pfVariant.ProductID,
		"variants":   variants,
		"name":       request.Name,
		"image":      request.Image,
	})

	if err != nil {
		log.Println(err)
		return nil, errors.New("error while calling printful api")
	}

	/*
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(body))
	*/
	response := CreateSyncProductResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error while decoding printful response")
	}

	if !response.Success {
		log.Println(response)
		return nil, errors.New("error while creating printful product")
	}

	log.Println("createProduct", response)

	products, err := createShopProduct(response.SyncProduct.ID)
	if err != nil {
		return nil, fmt.Errorf("error while creating shop product %w", err)
	}

	//return &variantResponse.Result.Variant, nil

	return products, nil
}

type CreateSyncProductResponse struct {
	Success     bool                `json:"success"`
	SyncProduct schemas.SyncProduct `json:"result"`
}

type GetSyncProductResponse struct {
	Success         bool                          `json:"success"`
	SyncProductInfo printfulModel.SyncProductInfo `json:"result"`
}

func createShopProduct(syncProductID int64) ([]*model.Product, error) {
	log.Println("creating product for id:", syncProductID)

	resp, err := fetchAPI("get-sync-product", 1, map[string]interface{}{
		"sync_product_id": syncProductID,
	})

	if err != nil {
		return nil, err
	}

	/*body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))*/
	response := GetSyncProductResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error while decoding printful response")
	}

	if !response.Success {
		log.Println(response)
		return nil, errors.New("error while creating printful product")
	}

	log.Println("createShopProduct", response)

	/*
		type SyncProductInfo struct {
			SyncProduct  SyncProduct   `json:"sync_product" bson:"sync_product"`
			SyncVariants []SyncVariant `json:"sync_variants" bson:"sync_variants"`
		}
	*/

	syncProduct := response.SyncProductInfo.SyncProduct
	syncVariants := response.SyncProductInfo.SyncVariants

	variantIDs := []string{}
	for _, syncVariant := range syncVariants {
		variantIDs = append(variantIDs, syncVariant.ExternalID)
	}

	shopProducts := []*model.Product{}
	for _, syncVariant := range syncVariants {
		//v = append(v, key)
		shopProduct, err := createShopProduct2(syncProduct, syncVariant, variantIDs)
		shopProducts = append(shopProducts, shopProduct)

		if err != nil {
			log.Println(err)
			return nil, errors.New("error while creating shop product")
		}

		log.Println(shopProduct)
		/*
			const shoProduct = await this.#createShopProduct2(syncProduct, syncVariants);
			productsIds.push(shoProduct.id);
			products.push(shoProduct);
		*/
	}
	/*
		if (productsIds.length > 1) {
			for (const productId of productsIds) {
				const updateOneResult = await this.#productsCollection.updateOne({ _id: productId }, { $set: { variantIds: productsIds }});
			}
		}
	*/

	// return the first product created
	//return products[0];

	return shopProducts, nil
}

func createShopProduct2(syncProduct schemas.SyncProduct, syncVariant schemas.SyncVariant, variantIDs []string) (*model.Product, error) {
	product := model.NewProduct()
	product.Name = syncVariant.Name
	product.ProductName = syncProduct.Name
	product.Currency = syncVariant.Currency
	product.ThumbnailURL = syncProduct.ThumbnailURL
	product.ExternalVariantID = syncVariant.ID
	product.Status = "completed"
	product.VariantIDs = variantIDs

	retailPrice, err := decimal.NewFromString(syncVariant.RetailPrice)
	if err != nil {
		return nil, err
	}
	product.RetailPrice = retailPrice

	id, err := primitive.ObjectIDFromHex(syncVariant.ExternalID)
	if err != nil {
		return nil, err
	}
	product.ID = id

	pfVariant, err := getPrintfulVariant(syncVariant.VariantID)
	if err != nil {
		return nil, err
	}

	log.Println(pfVariant)

	if pfVariant.ColorCode != "" {
		product.AddOption("color", "color", pfVariant.ColorCode)
	}
	if pfVariant.ColorCode2 != "" {
		product.AddOption("color2", "color", pfVariant.ColorCode2)
	}
	if pfVariant.Size != "" {
		product.AddOption("size", "size", pfVariant.Size)
	}
	if pfVariant.Image != "" {
		product.AddFile("product", pfVariant.Image)
	}

	pfProduct, err := getPrintfulProduct(pfVariant.ProductID)
	if err != nil {
		return nil, err
	}

	if pfProduct.Description != "" {
		product.Description = pfProduct.Description
	}

	for _, file := range syncVariant.Files {
		//v = append(v, key)
		product.AddFile(file.Type, file.URL)
	}

	err = mongo.UpdateProduct(&product)
	if err != nil {
		return nil, err
	}

	/*
		const shopProduct = new ShopProduct();
		const printfulProductReference = syncVariant.product;

		const printfulVariant = await this.#getPrintfulVariant(printfulProductReference.productId, printfulProductReference.variantId);
		if (!printfulVariant) {
			throw new Error(`Printful variant not found productId: ${printfulProductReference.productId} variantId: ${printfulProductReference.variantId}`);
		}

		const description = await this.#getPrintfulProductDescription(syncVariant?.product?.productId)
		if (description) {
			shopProduct.description = description;
		}

		const syncVariantFiles = syncVariant.files;
		if (syncVariantFiles) {
			for (const syncVariantFile of syncVariantFiles) {
				shopProduct.addFile(syncVariantFile.type, syncVariantFile.url);
			}
		}

		//console.log('createShopProduct2 printfulVariant', printfulVariant);
		const replaceOneResult = await this.#productsCollection.replaceOne({ _id: shopProduct.id }, shopProduct.toJSON());
		if (!replaceOneResult?.acknowledged) {
			winston.error('Error in #createShopProduct2 : replaceOne failed', { replaceOneResult: replaceOneResult, shopProduct: shopProduct.toJSON() });
			throw new Error('Error in #createShopProduct2 : replaceOne failed');
		}
		return shopProduct;
	*/
	/*
		if (productsIds.length > 1) {
			for (const productId of productsIds) {
				const updateOneResult = await this.#productsCollection.updateOne({ _id: productId }, { $set: { variantIds: productsIds }});
				//console.log(updateOneResult);
			}
		}
	*/

	return &product, nil
}

func createShopProducts(count int) ([]string, error) {
	ret := make([]string, 0, count)
	i := 0
	for i < count {
		product, err := mongo.CreateProduct()
		if err != nil {
			return nil, err
		}

		ret = append(ret, product.ID.Hex())

		i += 1
	}

	return ret, nil
}

func getPrintfulVariant(variantID int) (*printfulModel.Variant, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	/*u, err := url.JoinPath(printfulConfig.Endpoint, "/products/variant/", strconv.Itoa(int(variantID)))
	if err != nil {
		log.Println(err)
		return nil, errors.New("error while getting printful url")
	}

	log.Println(u)
	resp, err := http.Get(u)*/
	resp, err := fetchAPI("get-variant", 1, map[string]interface{}{
		"variant_id": variantID,
	})

	//body, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(body))

	if err != nil {
		log.Println(err)
		return nil, errors.New("error while calling printful api")
	}

	variantResponse := printfulModel.VariantResponse{}
	err = json.NewDecoder(resp.Body).Decode(&variantResponse)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error while decoding printful response")
	}

	if !variantResponse.Success {
		log.Println(variantResponse)
		return nil, errors.New("error while getting printful variant")
	}
	//log.Println("variantResponse", variantResponse)

	return &variantResponse.Result.Variant, nil
}

func getPrintfulProduct(productID int) (*printfulModel.Product, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	/*u, err := url.JoinPath(printfulConfig.Endpoint, "/product/", strconv.Itoa(int(productID)))
	if err != nil {
		log.Println(err)
		return nil, errors.New("error while getting printful url")
	}

	log.Println(u)
	resp, err := http.Get(u)*/
	resp, err := fetchAPI("get-product", 1, map[string]interface{}{
		"product_id": productID,
	})

	if err != nil {
		log.Println(err)
		return nil, errors.New("error while calling printful api")
	}

	productResponse := printfulModel.ProductResponse{}
	err = json.NewDecoder(resp.Body).Decode(&productResponse)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error while decoding printful response")
	}

	if !productResponse.Success {
		log.Println(productResponse)
		return nil, errors.New("error while getting printful variant")
	}

	return &productResponse.Result.Product, nil
}

func apiGetUserInfo(c *gin.Context, s sessions.Session) error {
	jsonSuccess(c, s.Get("user_infos"))
	return nil
}

func apiSetShippingAddress(c *gin.Context, s sessions.Session, params map[string]interface{}) error {

	log.Println(s)
	address := model.Address{}
	err := mapstructure.Decode(params["shipping_address"], &address)
	if err != nil {
		log.Println(err)
		return errors.New("error while reading params")
	}

	log.Println(address)
	orderID, ok := s.Get("order_id").(string)
	if !ok {
		return errors.New("error while retrieving order id")
	}
	order, err := mongo.FindOrder(orderID)
	if err != nil {
		log.Println(err)
		return errors.New("error while retrieving order")
	}

	order.ShippingAddress = address
	err = mongo.UpdateOrder(order)
	if err != nil {
		log.Println(err)
		return errors.New("error while updating order")
	}

	calculateShippingRatesRequest := printfulModel.CalculateShippingRatesRequest{Items: []printfulModel.ItemInfo{}}
	calculateShippingRatesRequest.Recipient.Address1 = order.ShippingAddress.Address1
	calculateShippingRatesRequest.Recipient.City = order.ShippingAddress.City
	calculateShippingRatesRequest.Recipient.CountryCode = order.ShippingAddress.CountryCode
	calculateShippingRatesRequest.Recipient.StateCode = order.ShippingAddress.StateCode
	calculateShippingRatesRequest.Recipient.ZIP = order.ShippingAddress.PostalCode

	for _, orderItem := range order.Items {
		itemInfo := printfulModel.ItemInfo{
			ExternalVariantID: orderItem.ProductID,
			Quantity:          int(orderItem.Quantity),
		}

		calculateShippingRatesRequest.Items = append(calculateShippingRatesRequest.Items, itemInfo)
	}

	resp, err := fetchAPI("calculate-shipping-rates", 1, calculateShippingRatesRequest)
	if err != nil {
		log.Println(err)
		return errors.New("error while calling printful api")
	}
	defer resp.Body.Close()

	response := calculateShippingRatesResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println(err)
		return errors.New("error while decoding printful response")
	}

	if !response.Success {
		log.Println(response)
		return errors.New("error while calculating shipping rates")
	}

	log.Println(order)
	order.ShippingInfos = response.ShippingInfos
	for _, shippingInfo := range order.ShippingInfos {
		order.ShippingMethod = shippingInfo.ID
		break
	}

	err = computeTaxRate(order)
	if err != nil {
		log.Println(err)
		return errors.New("error while computing shipping address")
	}

	err = mongo.UpdateOrder(order)
	if err != nil {
		log.Println(err)
		return errors.New("error while updating order")
	}

	jsonSuccess(c, map[string]interface{}{"order": order})
	return nil
}

func computeTaxRate(order *model.Order) error {
	calculateTaxRates := requestbodies.CalculateTaxRates{
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
	Success       bool                   `json:"success"`
	ShippingInfos []schemas.ShippingInfo `json:"result"`
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

	calculateShippingRatesRequest := printfulModel.CalculateShippingRatesRequest{Items: []printfulModel.ItemInfo{}}
	calculateShippingRatesRequest.Recipient.Address1 = order.ShippingAddress.Address1
	calculateShippingRatesRequest.Recipient.City = order.ShippingAddress.City
	calculateShippingRatesRequest.Recipient.CountryCode = order.ShippingAddress.CountryCode
	calculateShippingRatesRequest.Recipient.StateCode = order.ShippingAddress.StateCode
	calculateShippingRatesRequest.Recipient.ZIP = order.ShippingAddress.PostalCode

	for _, orderItem := range order.Items {
		itemInfo := printfulModel.ItemInfo{
			ExternalVariantID: orderItem.ProductID,
			Quantity:          int(orderItem.Quantity),
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

func apiSetShippingMethod(c *gin.Context, s sessions.Session, params map[string]interface{}) error {
	log.Println(s)

	method, ok := params["method"].(string)
	if !ok {
		return errors.New("error while getting shipping method")
	}

	log.Println(method)
	orderID, ok := s.Get("order_id").(string)
	if !ok {
		return errors.New("error while retrieving order id")
	}

	order, err := mongo.FindOrder(orderID)
	if err != nil {
		log.Println(err)
		return errors.New("error while retrieving order")
	}

	order.ShippingMethod = method
	err = mongo.UpdateOrder(order)
	if err != nil {
		log.Println(err)
		return errors.New("error while updating order")
	}

	err = createPrintfulOrder(*order)
	if err != nil {
		log.Println(err)
		return errors.New("error while creating printful order")
	}

	jsonSuccess(c, map[string]interface{}{"order": order})
	return nil
}

type createPrintfulOrderResponse struct {
	Success bool          `json:"success"`
	Order   schemas.Order `json:"result"`
}

func createPrintfulOrder(order model.Order) error {
	printfulOrder := schemas.NewOrder()
	printfulOrder.Recipient.Address1 = order.ShippingAddress.Address1
	printfulOrder.Recipient.City = order.ShippingAddress.City
	printfulOrder.Recipient.CountryCode = order.ShippingAddress.CountryCode
	printfulOrder.Recipient.StateCode = order.ShippingAddress.StateCode
	printfulOrder.Recipient.ZIP = order.ShippingAddress.PostalCode

	log.Println(printfulOrder)
	/*
		calculateShippingRatesRequest.Recipient.Address1 = order.ShippingAddress.Address1
		calculateShippingRatesRequest.Recipient.City = order.ShippingAddress.City
		calculateShippingRatesRequest.Recipient.CountryCode = order.ShippingAddress.CountryCode
		calculateShippingRatesRequest.Recipient.StateCode = order.ShippingAddress.StateCode
		calculateShippingRatesRequest.Recipient.ZIP = order.ShippingAddress.PostalCode
	*/

	for _, orderItem := range order.Items {
		log.Println("**********************", orderItem)
		item := schemas.Item{
			ExternalVariantID: orderItem.ProductID,
			Quantity:          int(orderItem.Quantity),
			RetailPrice:       orderItem.RetailPrice.String(),
		}
		log.Println("AAAAAAAAAAAAAAAAAAAAAA", orderItem.RetailPrice.String())
		printfulOrder.Items = append(printfulOrder.Items, item)
	}

	resp, err := fetchAPI("create-order", 1, map[string]interface{}{
		"order": printfulOrder,
	})
	if err != nil {
		log.Println(err)
		return errors.New("error while calling printful api")
	}
	defer resp.Body.Close()

	response := createPrintfulOrderResponse{}
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

/*
roundPrice(currency, price) {
	let digits = CURRENCIES_DIGITS[currency] ?? 2;
	return Number(Number.parseFloat(price).toFixed(digits));
}
*/

func apiCreatePaypalOrder(c *gin.Context, s sessions.Session) error {
	//log.Println(s)

	orderID, ok := s.Get("order_id").(string)
	if !ok {
		return errors.New("error while retrieving order id")
	}

	order, err := mongo.FindOrder(orderID)
	if err != nil {
		log.Println(err)
		return errors.New("error while retrieving order")
	}

	fmt.Println(order)

	client, err := paypal.NewClient(paypalConfig.ClientID, paypalConfig.ClientSecret, paypal.APIBaseSandBox)
	if err != nil {
		log.Println(err)
		return errors.New("error while creating paypal client")
	}

	paypalOrder, err := client.CreateOrder(
		context.Background(),
		paypal.OrderIntentCapture,
		[]paypal.PurchaseUnitRequest{
			{
				Amount: &paypal.PurchaseUnitAmount{
					Value:    order.GetTotalPrice().String(),
					Currency: order.Currency,
					Breakdown: &paypal.PurchaseUnitAmountBreakdown{
						ItemTotal: &paypal.Money{
							Currency: order.Currency,
							Value:    order.GetItemsPrice().String(),
						},
						Shipping: &paypal.Money{
							Currency: order.Currency,
							Value:    order.GetShippingPrice().String(),
						},
						TaxTotal: &paypal.Money{
							Currency: order.Currency,
							Value:    order.GetTaxPrice().String(),
						},
					},
					/*
						amount: {
							currency_code: currency,
							value: roundPrice(currency, order.totalPrice),
							breakdown: {
							}
						},
					*/
				},
				CustomID: order.ID.Hex(),
				Shipping: &paypal.ShippingDetail{
					Name: &paypal.Name{
						FullName: order.ShippingAddress.GetFullName(),
					},
					Address: &paypal.ShippingDetailAddressPortable{
						AddressLine1: order.ShippingAddress.Address1,
						AddressLine2: order.ShippingAddress.Address2,
						AdminArea1:   order.ShippingAddress.StateCode,
						AdminArea2:   order.ShippingAddress.City,
						PostalCode:   order.ShippingAddress.PostalCode,
						CountryCode:  order.ShippingAddress.CountryCode,
					},
				},
			},
		},
		&paypal.CreateOrderPayer{},
		&paypal.ApplicationContext{
			ShippingPreference: paypal.ShippingPreferenceSetProvidedAddress,
		},
	)

	if err != nil {
		log.Println(err)
		return errors.New("error while creating paypal order")
	}

	log.Println("Got paypal order:", paypalOrder)

	order.PaypalOrderID = paypalOrder.ID
	err = mongo.UpdateOrder(order)
	if err != nil {
		log.Println(err)
		return errors.New("error while updating order")
	}

	jsonSuccess(c, map[string]interface{}{"paypal_order_id": paypalOrder.ID})
	return nil
}

func apiCapturePaypalOrder(c *gin.Context, s sessions.Session, params map[string]interface{}) error {
	if params == nil {
		return errors.New("no params provided")
	}

	var id interface{}
	var ok bool
	if id, ok = params["paypal_order_id"]; !ok {
		return errors.New("missing param paypal_order_id")
	}

	orderId := id.(string)

	if len(orderId) > 36 {
		return errors.New("paypal order id is too long")
	}
	if !IsAlphaNumeric(orderId) {
		return errors.New("paypal order id has a wrong format " + orderId)
	}

	client, err := paypal.NewClient(paypalConfig.ClientID, paypalConfig.ClientSecret, paypal.APIBaseSandBox)
	if err != nil {
		log.Println(err)
		return errors.New("error while creating paypal client")
	}

	paypalOrder, err := client.GetOrder(
		context.Background(),
		orderId,
	)

	if err != nil {
		log.Println(err)
		return errors.New("error while retrieving paypal order")
	}

	if paypalOrder.Status != "APPROVED" {
		return errors.New("paypal order is not approved")
	}

	order, err := mongo.FindOrderByPaypalID(orderId)
	if err != nil {
		log.Println(err)
		return errors.New("error while retrieving order")
	}

	order.Status = "approved"
	err = mongo.UpdateOrder(order)
	if err != nil {
		log.Println(err)
		return errors.New("error while updating order")
	}

	cart := s.Get("cart").(model.Cart)
	cart.Clear()

	jsonSuccess(c, map[string]interface{}{"order": order})
	return nil
}
