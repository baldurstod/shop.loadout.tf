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
	"io/ioutil"
	_ "strconv"
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

	p, err := mongo.GetProduct(params["product_id"].(string))

	if err != nil {
		log.Println(err)
		return errors.New("Error while getting products")
	}

	jsonSuccess(w, r, map[string]interface{}{"product": p})
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
		orderItem.ThumbnailUrl = p.ThumbnailUrl
		orderItem.Quantity = quantity
		orderItem.RetailPrice = p.RetailPrice

		order.Items = append(order.Items, orderItem)
	}

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
		"variants": variants,
		"name":     request.Name,
		"image":    request.Image,
	})

	if err != nil {
		log.Println(err)
		return errors.New("Error while calling printful api")
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))

	return nil
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
