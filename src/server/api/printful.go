package api

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	printfulModel "github.com/baldurstod/printful-api-model"
	"log"
	"net/http"
	"net/url"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/model"
	"shop.loadout.tf/src/server/model/requests"
	"shop.loadout.tf/src/server/mongo"
	//"shop.loadout.tf/src/server/sessions"
	"bytes"
	"github.com/gorilla/sessions"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "io/ioutil"
	"strconv"
	"time"
)

var printfulConfig config.Printful
var printfulURL string

func SetPrintfulConfig(config config.Printful) {
	printfulConfig = config
	log.Println(config)
	var err error
	printfulURL, err = url.JoinPath(printfulConfig.Endpoint, "/api")
	if err != nil {
		panic("Error while getting printful url")
	}
}

func fetchAPI(action string, version int, params interface{}) (*http.Response, error) {

	body := map[string]interface{}{
		"action":  action,
		"version": version,
		"params":  params,
	}

	requestBody, err := json.Marshal(body)
	log.Println(string(requestBody))
	res, err := http.Post(printfulURL, "application/json", bytes.NewBuffer(requestBody))

	return res, err
}

func getCountries(w http.ResponseWriter, r *http.Request) error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	/*
		u, err := url.JoinPath(printfulConfig.Endpoint, "/countries")
		if err != nil {
			return errors.New("Error while getting printful url")
		}

		resp, err := http.Get(u)*/
	resp, err := fetchAPI("get-countries", 1, nil)
	//body, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(body))
	log.Println(resp)

	if err != nil {
		log.Println(err)
		return errors.New("Error while calling printful api")
	}
	defer resp.Body.Close()

	countriesResponse := printfulModel.CountriesResponse{}
	err = json.NewDecoder(resp.Body).Decode(&countriesResponse)
	if err != nil {
		log.Println(err)
		return errors.New("Error while decoding printful response")
	}

	jsonSuccess(w, r, map[string]interface{}{"countries": countriesResponse.Countries})

	return nil
}

func getCurrency(w http.ResponseWriter, r *http.Request, s *sessions.Session) error {
	jsonSuccess(w, r, s.Values["currency"])
	return nil
}

func getFavorites(w http.ResponseWriter, r *http.Request, s *sessions.Session) error {
	favorites := s.Values["favorites"].(map[string]interface{})

	v := make([]string, 0, len(favorites))

	for key := range favorites {
		v = append(v, key)
	}

	jsonSuccess(w, r, v)
	return nil
}

func getProduct(w http.ResponseWriter, r *http.Request, params map[string]interface{}) error {
	if params == nil {
		return errors.New("No params provided")
	}

	product, err := mongo.FindProduct(params["product_id"].(string))

	if err != nil {
		log.Println(err)
		return errors.New("Error while getting products")
	}

	for _, variantID := range product.VariantIDs {
		//variants[variantID] = struct{}{}
		p, err := mongo.FindProduct(variantID)

		if err == nil {
			product.AddVariant(model.NewVariant(p))
		}
	}

	jsonSuccess(w, r, map[string]interface{}{"product": product})
	return nil
}

func getProducts(w http.ResponseWriter, r *http.Request, s *sessions.Session) error {
	p, err := mongo.GetProducts()

	if err != nil {
		log.Println(err)
		return errors.New("Error while getting products")
	}

	jsonSuccess(w, r, p)
	return nil
}

func sendContact(w http.ResponseWriter, r *http.Request, params map[string]interface{}) error {
	if params == nil {
		return errors.New("No params provided")
	}

	id, err := mongo.SendContact(params)

	if err != nil {
		log.Println(err)
		return errors.New("Error while sending contact")
	}

	jsonSuccess(w, r, id)
	return nil
}

func setFavorite(w http.ResponseWriter, r *http.Request, s *sessions.Session, params map[string]interface{}) error {
	if params == nil {
		return errors.New("No params provided")
	}

	pID, ok := params["product_id"]
	isFavorite, ok2 := params["is_favorite"]

	if !ok || !ok2 {
		return errors.New("Missing params")
	}

	favorites := s.Values["favorites"].(map[string]interface{})

	productId := pID.(string)
	if isFavorite.(bool) {
		favorites[productId] = struct{}{}
	} else {
		delete(favorites, productId)
	}

	log.Println(favorites)

	jsonSuccess(w, r, nil)
	return nil
}

func addProduct(w http.ResponseWriter, r *http.Request, s *sessions.Session, params map[string]interface{}) error {
	if params == nil {
		return errors.New("No params provided")
	}

	pID, ok := params["product_id"]
	quantity, ok2 := params["quantity"]

	if !ok || !ok2 {
		return errors.New("Missing params")
	}

	cart := s.Values["cart"].(model.Cart)

	cart.AddQuantity(pID.(string), uint(quantity.(float64)))

	jsonSuccess(w, r, map[string]interface{}{"cart": cart})
	return nil
}

func setProductQuantity(w http.ResponseWriter, r *http.Request, s *sessions.Session, params map[string]interface{}) error {
	if params == nil {
		return errors.New("No params provided")
	}

	pID, ok := params["product_id"]
	quantity, ok2 := params["quantity"]

	if !ok || !ok2 {
		return errors.New("Missing params")
	}

	cart := s.Values["cart"].(model.Cart)

	cart.SetQuantity(pID.(string), uint(quantity.(float64)))

	jsonSuccess(w, r, map[string]interface{}{"cart": cart})
	return nil
}

func getCart(w http.ResponseWriter, r *http.Request, s *sessions.Session) error {
	cart := s.Values["cart"].(model.Cart)

	jsonSuccess(w, r, map[string]interface{}{"cart": cart})
	return nil
}

func initCheckout(w http.ResponseWriter, r *http.Request, s *sessions.Session, params map[string]interface{}) error {
	cart := s.Values["cart"].(model.Cart)

	order, err := mongo.CreateOrder()
	if err != nil {
		log.Println(err)
		return errors.New("Error while creating order")
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
		return errors.New("Error while adding items to order")
	}

	order.Currency = cart.Currency
	now := time.Now().Unix()
	order.DateCreated = now
	order.DateUpdated = now
	order.Status = "created"

	err = mongo.UpdateOrder(order)
	if err != nil {
		log.Println(err)
		return errors.New("Error while updating order")
	}

	log.Println(order)
	s.Values["order_id"] = order.ID.Hex();
	log.Println(s)
	saveSession(w, r, s)
	jsonSuccess(w, r, map[string]interface{}{"order": order})

	return nil
}

func initCheckoutItems(cart *model.Cart, order *model.Order) error {
	log.Println(cart.Items)
	for productID, quantity := range cart.Items {
		p, err := mongo.GetProduct(productID)
		if err != nil {
			log.Println(err)
			return errors.New("Error during order initialization")
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

func apiCreateProduct(w http.ResponseWriter, r *http.Request, s *sessions.Session, params map[string]interface{}) error {
	if params == nil {
		return errors.New("No params provided")
	}
	//log.Println(params)
	//createProduct := params["product"].(requests.CreateProductRequest)

	createProductRequest := requests.CreateProductRequest{}
	err := mapstructure.Decode(params["product"], &createProductRequest)
	if err != nil {
		log.Println(err)
		return errors.New("Error while reading params")
	}

	err = createProductRequest.CheckParams()
	if err != nil {
		log.Println(err)
		return errors.New("Invalid params")
	}

	log.Println(createProductRequest.Name, createProductRequest.Type, createProductRequest.VariantID)
	createProduct(&createProductRequest)

	return errors.New("Error while creating product")
}

func createProduct(request *requests.CreateProductRequest) error {
	pfVariant, err := getPrintfulVariant(request.VariantID)
	if err != nil {
		log.Println(err)
		return errors.New("Variant not found")
	}

	log.Println(pfVariant)
	pfProduct, err := getPrintfulProduct(pfVariant.ProductID)
	if err != nil {
		log.Println(err)
		return errors.New("Product not found")
	}

	log.Println(pfProduct)

	resp, err := fetchAPI("get-similar-variants", 1, map[string]interface{}{
		"variant_id": pfVariant.ID,
	})

	if err != nil {
		log.Println(err)
		return errors.New("Error while calling printful api")
	}

	similarVariantsResponse := printfulModel.SimilarVariantsResponse{}
	err = json.NewDecoder(resp.Body).Decode(&similarVariantsResponse)
	if err != nil {
		log.Println(err)
		return errors.New("Error while decoding printful response")
	}

	if !similarVariantsResponse.Success {
		log.Println(similarVariantsResponse)
		return errors.New("Error while getting printful variant")
	}

	log.Println(similarVariantsResponse)

	variantCount := len(similarVariantsResponse.SimilarVariants)
	ids, err := createShopProducts(variantCount)
	if err != nil {
		log.Println(err)
		return errors.New("Error while creating products")
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
		return errors.New("Error while calling printful api")
	}

	/*
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(body))
	*/
	response := CreateSyncProductResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println(err)
		return errors.New("Error while decoding printful response")
	}

	if !response.Success {
		log.Println(response)
		return errors.New("Error while creating printful product")
	}

	log.Println("createProduct", response)

	createShopProduct(response.SyncProduct.ID)

	//return &variantResponse.Result.Variant, nil

	return nil
}

type CreateSyncProductResponse struct {
	Success     bool                      `json:"success"`
	SyncProduct printfulModel.SyncProduct `json:"result"`
}

type GetSyncProductResponse struct {
	Success         bool                          `json:"success"`
	SyncProductInfo printfulModel.SyncProductInfo `json:"result"`
}

func createShopProduct(syncProductID int64) error {
	log.Println("creating product for id:", syncProductID)

	resp, err := fetchAPI("get-sync-product", 1, map[string]interface{}{
		"sync_product_id": syncProductID,
	})

	if err != nil {
		return err
	}

	/*body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))*/
	response := GetSyncProductResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println(err)
		return errors.New("Error while decoding printful response")
	}

	if !response.Success {
		log.Println(response)
		return errors.New("Error while creating printful product")
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

	for _, syncVariant := range syncVariants {
		//v = append(v, key)
		shoProduct, err := createShopProduct2(syncProduct, syncVariant, variantIDs)

		if err != nil {
			log.Println(err)
			return errors.New("Error while creating shop product")
		}

		log.Println(shoProduct)
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

	return nil
}

func createShopProduct2(syncProduct printfulModel.SyncProduct, syncVariant printfulModel.SyncVariant, variantIDs []string) (*model.Product, error) {
	product := model.NewProduct()
	product.Name = syncVariant.Name
	product.ProductName = syncProduct.Name
	product.Currency = syncVariant.Currency
	product.ThumbnailURL = syncProduct.ThumbnailURL
	product.ExternalVariantID = syncVariant.ID
	product.Status = "completed"
	product.VariantIDs = variantIDs

	retailPrice, err := strconv.ParseFloat(syncVariant.RetailPrice, 32)
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
		return nil, errors.New("Error while getting printful url")
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
		return nil, errors.New("Error while calling printful api")
	}

	variantResponse := printfulModel.VariantResponse{}
	err = json.NewDecoder(resp.Body).Decode(&variantResponse)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error while decoding printful response")
	}

	if !variantResponse.Success {
		log.Println(variantResponse)
		return nil, errors.New("Error while getting printful variant")
	}
	//log.Println("variantResponse", variantResponse)

	return &variantResponse.Result.Variant, nil
}

func getPrintfulProduct(productID int) (*printfulModel.Product, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	/*u, err := url.JoinPath(printfulConfig.Endpoint, "/product/", strconv.Itoa(int(productID)))
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error while getting printful url")
	}

	log.Println(u)
	resp, err := http.Get(u)*/
	resp, err := fetchAPI("get-product", 1, map[string]interface{}{
		"product_id": productID,
	})

	if err != nil {
		log.Println(err)
		return nil, errors.New("Error while calling printful api")
	}

	productResponse := printfulModel.ProductResponse{}
	err = json.NewDecoder(resp.Body).Decode(&productResponse)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error while decoding printful response")
	}

	if !productResponse.Success {
		log.Println(productResponse)
		return nil, errors.New("Error while getting printful variant")
	}

	return &productResponse.Result.Product, nil
}

func apiGetUserInfo(w http.ResponseWriter, r *http.Request, s *sessions.Session, params map[string]interface{}) error {
	jsonSuccess(w, r, s.Values["user_infos"])
	return nil
}

func apiSetShippingAddress(w http.ResponseWriter, r *http.Request, s *sessions.Session, params map[string]interface{}) error {

	log.Println(s)
	address := model.Address{}
	err := mapstructure.Decode(params["shipping_address"], &address)
	if err != nil {
		log.Println(err)
		return errors.New("Error while reading params")
	}

	log.Println(address)
	orderID, ok := s.Values["order_id"].(string)
	if !ok {
		return errors.New("Error while retrieving order id")
	}
	order, err := mongo.FindOrder(orderID)
	if err != nil {
		log.Println(err)
		return errors.New("Error while retrieving order")
	}

	order.ShippingAddress = address
	err = mongo.UpdateOrder(order)
	if err != nil {
		log.Println(err)
		return errors.New("Error while updating order")
	}

	log.Println(order)
	jsonSuccess(w, r, map[string]interface{}{"order": order})
	return nil
}
