package api

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"shop.loadout.tf/src/server/model"
	"shop.loadout.tf/src/server/mongo"
	"shop.loadout.tf/src/server/printful"
)

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

	jsonSuccess(c, map[string]interface{}{"favorites": v})
	return nil
}

type ProductPrice struct {
	Currency string            `json:"currency" bson:"currency" mapstructure:"currency"`
	Prices   map[string]string `json:"prices" bson:"prices" mapstructure:"prices"`
}

func NewProductPrice(currency string) *ProductPrice {
	return &ProductPrice{
		Currency: currency,
		Prices:   map[string]string{},
	}
}

func getProduct(c *gin.Context, s sessions.Session, params map[string]interface{}) error {
	if params == nil {
		return errors.New("no params provided")
	}

	currency := s.Get("currency").(string)
	prices := NewProductPrice(currency)

	productID, ok := params["product_id"].(string)
	if !ok {
		return errors.New("invalid product id")
	}
	product, err := mongo.FindProduct(productID)

	if err != nil {
		log.Println(err)
		return errors.New("error while getting product")
	}
	prices.Prices[product.ID] = ""

	for _, variantID := range product.VariantIDs {
		//variants[variantID] = struct{}{}
		p, err := mongo.FindProduct(variantID)

		if err == nil {
			product.AddVariant(model.NewVariant(p))
		}
	}

	for id := range prices.Prices {
		price, err := mongo.GetRetailPrice(id, currency)
		if err != nil {
			log.Println(err)
		}
		prices.Prices[id] = price.RetailPrice.String()
	}

	jsonSuccess(c, map[string]any{
		"product": product,
		"prices":  prices,
	})
	return nil
}

func getProducts(c *gin.Context, s sessions.Session) error {
	p, err := mongo.GetProducts()

	currency := s.Get("currency").(string)
	prices := NewProductPrice(currency)
	for _, p2 := range p {
		prices.Prices[p2.ID] = ""
		for _, id := range p2.VariantIDs {
			prices.Prices[id] = ""
		}
	}

	for id := range prices.Prices {
		price, err := mongo.GetRetailPrice(id, currency)
		if err != nil {
			log.Println(err)
		}
		prices.Prices[id] = price.RetailPrice.String()
	}

	if err != nil {
		log.Println(err)
		return errors.New("error while getting products")
	}

	jsonSuccess(c, map[string]any{
		"products": p,
		"prices":   prices,
	})
	return nil
}

func sendContact(c *gin.Context, params map[string]interface{}) error {
	if params == nil {
		return errors.New("no params provided")
	}

	subject, ok := params["subject"].(string)
	if !ok {
		return errors.New("invalid subject")
	}

	email, ok := params["email"].(string)
	if !ok {
		return errors.New("invalid email")
	}

	content, ok := params["content"].(string)
	if !ok {
		return errors.New("invalid content")
	}

	id, err := mongo.SendContact(subject, email, content)

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

	productId, ok := params["product_id"].(string)
	isFavorite, ok2 := params["is_favorite"].(bool)

	if !ok || !ok2 {
		return errors.New("missing params")
	}

	favorites := s.Get("favorites").(map[string]interface{})

	if isFavorite {
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

	productId, ok := params["product_id"].(string)
	quantity, ok2 := params["quantity"].(float64)

	if !ok || !ok2 {
		return errors.New("missing params")
	}

	cart := s.Get("cart").(model.Cart)

	cart.AddQuantity(productId, uint(quantity))
	s.Delete("order_id")

	jsonSuccess(c, map[string]interface{}{"cart": cart})
	return nil
}

func setProductQuantity(c *gin.Context, s sessions.Session, params map[string]interface{}) error {
	if params == nil {
		return errors.New("no params provided")
	}

	productId, ok := params["product_id"].(string)
	quantity, ok2 := params["quantity"].(float64)

	if !ok || !ok2 {
		return errors.New("missing params")
	}

	cart := s.Get("cart").(model.Cart)

	cart.SetQuantity(productId, uint(quantity))
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

	order.Currency = cart.Currency
	err = initCheckoutItems(&cart, order)
	if err != nil {
		log.Println(err)
		return errors.New("error while adding items to order")
	}

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
	s.Set("order_id", order.ID)
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

		price, err := mongo.GetRetailPrice(productID, order.Currency)
		if err != nil {
			log.Println(err)
			return errors.New("error during order initialization")
		}

		orderItem := model.OrderItem{}
		orderItem.ProductID = p.ID
		orderItem.Product = *p
		orderItem.Name = p.Name
		orderItem.ThumbnailURL = p.ThumbnailURL
		orderItem.Quantity = quantity
		orderItem.RetailPrice = price.RetailPrice

		order.Items = append(order.Items, orderItem)
	}

	log.Println("-----------------", order)

	return nil
}

func apiGetUserInfo(c *gin.Context, s sessions.Session) error {
	jsonSuccess(c, s.Get("user_infos"))
	return nil
}

func apiSetShippingAddress(c *gin.Context, s sessions.Session, params map[string]any) error {
	log.Println(s)
	shippingAddress := model.Address{}
	billingAddress := model.Address{}

	if params["shipping_address"] == nil {
		return errors.New("missing param shipping_address")
	}

	err := mapstructure.Decode(params["shipping_address"], &shippingAddress)
	if err != nil {
		log.Println(err)
		return errors.New("error while reading params")
	}

	if err := checkAddress(&shippingAddress); err != nil {
		log.Println(err)
		return fmt.Errorf("incomplete shipping adress: %v", err)
	}

	sameBillingAddress, ok := params["same_billing_address"].(bool)
	if !ok {
		return errors.New("error while reading param same_billing_address")
	}

	if !sameBillingAddress {
		if params["billing_address"] == nil {
			return errors.New("missing param billing_address")
		}

		err := mapstructure.Decode(params["billing_address"], &billingAddress)
		if err != nil {
			log.Println(err)
			return errors.New("error while reading param billing_address")
		}

		if err := checkAddress(&billingAddress); err != nil {
			log.Println(err)
			return fmt.Errorf("incomplete billing adress: %v", err)
		}
	}

	log.Println(shippingAddress)
	orderID, ok := s.Get("order_id").(string)
	if !ok {
		return errors.New("error while retrieving order id")
	}
	order, err := mongo.FindOrder(orderID)
	if err != nil {
		log.Println(err)
		return errors.New("error while retrieving order")
	}

	if order.Status == "approved" {
		return fmt.Errorf("error %s is already approved", orderID)
	}

	order.ShippingAddress = shippingAddress
	order.SameBillingAddress = sameBillingAddress
	order.BillingAddress = billingAddress

	err = mongo.UpdateOrder(order)
	if err != nil {
		log.Println(err)
		return errors.New("error while updating order")
	}

	recipient := printfulmodel.ShippingRatesAddress{
		Address1:    order.ShippingAddress.Address1,
		City:        order.ShippingAddress.City,
		CountryCode: order.ShippingAddress.CountryCode,
		StateCode:   order.ShippingAddress.StateCode,
		ZIP:         order.ShippingAddress.PostalCode,
	}

	items := []printfulmodel.CatalogOrWarehouseShippingRateItem{}

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

		items = append(items, itemInfo)
	}

	shippingInfos, err := printful.CalculateShippingRates(recipient, items, "", "") /*TODO: add currency, locale*/
	if err != nil {
		log.Println(err)
		return errors.New("error while computing shipping rates in apiSetShippingAddress")
	}

	log.Println(order)
	order.ShippingInfos = shippingInfos
	for _, shippingInfo := range order.ShippingInfos {
		order.ShippingMethod = shippingInfo.Shipping
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

func checkAddress(address *model.Address) error {
	if address.FirstName == "" {
		return errors.New("first name is missing")
	}

	if address.LastName == "" {
		return errors.New("last name is missing")
	}

	if address.Address1 == "" {
		return errors.New("address line 1 is missing")
	}

	if address.City == "" {
		return errors.New("city is missing")
	}

	if address.CountryCode == "" {
		return errors.New("country code is missing")
	}

	if address.PostalCode == "" {
		return errors.New("postal code is missing")
	}

	if address.Phone == "" {
		return errors.New("phone number is missing")
	}

	if address.Email == "" {
		return errors.New("email is missing")
	}

	return nil
}

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

	if order.Status == "approved" {
		return fmt.Errorf("error %s is already approved", orderID)
	}

	order.ShippingMethod = method
	err = mongo.UpdateOrder(order)
	if err != nil {
		log.Println(err)
		return errors.New("error while updating order")
	}

	/*
		err = createPrintfulOrder(*order)
		if err != nil {
			log.Println(err)
			return errors.New("error while creating printful order")
		}
	*/

	jsonSuccess(c, map[string]interface{}{"order": order})
	return nil
}

func apiGetOrder(c *gin.Context, params map[string]interface{}) error {
	if params == nil {
		return errors.New("no params provided")
	}

	orderID, ok := params["order_id"].(string)
	if !ok {
		return errors.New("invalid order id")
	}

	order, err := mongo.FindOrder(orderID)
	if err != nil {
		log.Println(err)
		return errors.New("error while getting order")
	}

	jsonSuccess(c, map[string]interface{}{"order": order})
	return nil
}
