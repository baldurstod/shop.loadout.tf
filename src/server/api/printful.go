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
	"github.com/gorilla/sessions"
	"github.com/mitchellh/mapstructure"
	"time"
)

var printfulConfig config.Printful

func SetPrintfulConfig(config config.Printful) {
	printfulConfig = config
	log.Println(config)
}

func getCountries(w http.ResponseWriter, r *http.Request) error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	u, err := url.JoinPath(printfulConfig.Endpoint, "/countries")
	if err != nil {
		return errors.New("Error while getting printful url")
	}

	resp, err := http.Get(u)

	if err != nil {
		log.Println(err)
		return errors.New("Error while calling printful api")
	}

	countriesResponse := printfulModel.CountriesResponse{}
	err = json.NewDecoder(resp.Body).Decode(&countriesResponse)
	if err != nil {
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

func createProduct(w http.ResponseWriter, r *http.Request, s *sessions.Session, params map[string]interface{}) error {
	if params == nil {
		return errors.New("No params provided")
	}
	//log.Println(params)
	//createProduct := params["product"].(requests.CreateProductRequest)

	createProduct := requests.CreateProductRequest{}
	err := mapstructure.Decode(params["product"], &createProduct)
	if err != nil {
		errors.New("Error while reading params")
	}

	log.Println(createProduct.Name, createProduct.Type, createProduct.VariantID)

	return errors.New("Error while creating product")
}
