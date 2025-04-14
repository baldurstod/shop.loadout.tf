package api

import (
	"errors"
	"fmt"
	"log"
	"net/mail"
	"strconv"
	"time"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"shop.loadout.tf/src/server/constants"
	"shop.loadout.tf/src/server/databases"
	"shop.loadout.tf/src/server/logger"
	"shop.loadout.tf/src/server/model"
	"shop.loadout.tf/src/server/printful"
)

func apiGetCurrency(c *gin.Context, s sessions.Session) error {
	currency, ok := s.Get("currency").(string)
	if !ok {
		currency = constants.DEFAULT_CURRENCY
	}

	jsonSuccess(c, map[string]any{"currency": currency})
	return nil
}

func apiGetFavorites(c *gin.Context, s sessions.Session) error {
	favorites, ok := s.Get("favorites").(map[string]any)
	if !ok {
		favorites = make(map[string]any)
	}

	v := make([]string, 0, len(favorites))

	for key := range favorites {
		v = append(v, key)
	}

	jsonSuccess(c, map[string]any{"favorites": v})
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

func apiGetProduct(c *gin.Context, s sessions.Session, params map[string]any) error {
	if params == nil {
		return errors.New("no params provided")
	}

	currency, ok := s.Get("currency").(string)
	if !ok {
		currency = constants.DEFAULT_CURRENCY
	}
	prices := NewProductPrice(currency)

	productID, ok := params["product_id"].(string)
	if !ok {
		return errors.New("invalid product id")
	}

	product, err := databases.FindProduct(productID)
	if err != nil {
		logger.Log(c, err)
		return errors.New("error while getting product")
	}
	prices.Prices[product.ID] = ""

	for _, variantID := range product.VariantIDs {
		//variants[variantID] = struct{}{}
		p, err := databases.FindProduct(variantID)

		if err == nil {
			product.AddVariant(model.NewVariant(p))
		}
	}

	for id := range prices.Prices {
		price, err := databases.GetRetailPrice(id, currency)
		if err != nil {
			logger.Log(c, err)
		}
		prices.Prices[id] = price.RetailPrice.String()
	}

	jsonSuccess(c, map[string]any{
		"product": product,
		"prices":  prices,
	})
	return nil
}

func apiGetProducts(c *gin.Context, s sessions.Session) error {
	p, err := databases.GetProducts()

	currency, ok := s.Get("currency").(string)
	if !ok {
		currency = constants.DEFAULT_CURRENCY
	}
	prices := NewProductPrice(currency)
	for _, p2 := range p {
		prices.Prices[p2.ID] = ""
		for _, id := range p2.VariantIDs {
			prices.Prices[id] = ""
		}
	}

	for id := range prices.Prices {
		price, err := databases.GetRetailPrice(id, currency)
		if err != nil {
			logger.Log(c, err)
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

func validEmail(email string) bool {
	emailAddress, err := mail.ParseAddress(email)
	return err == nil && emailAddress.Address == email
}

func apiSendMessage(c *gin.Context, params map[string]any) error {
	if params == nil {
		return errors.New("no params provided")
	}

	subject, ok := params["subject"].(string)
	if !ok || subject == "" || len(subject) < 5 {
		return errors.New("invalid subject")
	}

	email, ok := params["email"].(string)
	if !ok || email == "" || !validEmail(email) {
		return errors.New("invalid email")
	}

	content, ok := params["content"].(string)
	if !ok || content == "" || len(content) < 10 {
		return errors.New("invalid content")
	}

	id, err := databases.SendContact(subject, email, content)

	if err != nil {
		logger.Log(c, err)
		return errors.New("error while sending message")
	}

	jsonSuccess(c, map[string]any{"message_id": id})
	return nil
}

func apiSetFavorite(c *gin.Context, s sessions.Session, params map[string]any) error {
	if params == nil {
		return errors.New("no params provided")
	}

	productId, ok := params["product_id"].(string)
	if !ok {
		return errors.New("missing params product_id")
	}

	isFavorite, ok := params["is_favorite"].(bool)
	if !ok {
		return errors.New("missing params is_favorite")
	}

	favorites, ok := s.Get("favorites").(map[string]any)
	if !ok {
		return errors.New("favorites not found")
	}

	if isFavorite {
		favorites[productId] = struct{}{}
	} else {
		delete(favorites, productId)
	}

	log.Println(favorites)

	jsonSuccess(c, nil)
	return nil
}

func apiAddProduct(c *gin.Context, s sessions.Session, params map[string]any) error {
	if params == nil {
		return errors.New("no params provided")
	}

	productId, ok := params["product_id"].(string)
	if !ok {
		return errors.New("missing params product_id")
	}

	quantity, ok := params["quantity"].(float64)
	if !ok {
		return errors.New("missing params quantity")
	}

	cart, ok := s.Get("cart").(model.Cart)
	if !ok {
		err := errors.New("cart not found")
		logger.Log(c, err)
		return err
	}

	cart.AddQuantity(productId, uint(quantity))
	s.Delete("order_id")

	jsonSuccess(c, map[string]any{"cart": cart})
	return nil
}

func apiSetProductQuantity(c *gin.Context, s sessions.Session, params map[string]any) error {
	if params == nil {
		return errors.New("no params provided")
	}

	productId, ok := params["product_id"].(string)
	if !ok {
		return errors.New("missing params product_id")
	}

	quantity, ok := params["quantity"].(float64)
	if !ok {
		return errors.New("missing params quantity")
	}

	cart, ok := s.Get("cart").(model.Cart)
	if !ok {
		err := errors.New("cart not found")
		logger.Log(c, err)
		return err
	}

	cart.SetQuantity(productId, uint(quantity))
	s.Delete("order_id")

	jsonSuccess(c, map[string]any{"cart": cart})
	return nil
}

func apiGetCart(c *gin.Context, s sessions.Session) error {
	cart, ok := s.Get("cart").(model.Cart)
	if !ok {
		cart = model.NewCart()
	}

	jsonSuccess(c, map[string]any{"cart": cart})
	return nil
}

func apiInitCheckout(c *gin.Context, s sessions.Session) error {
	cart, ok := s.Get("cart").(model.Cart)
	if !ok {
		err := errors.New("cart not found")
		logger.Log(c, err)
		return err
	}

	order, err := databases.CreateOrder()
	if err != nil {
		logger.Log(c, err)
		return errors.New("error while creating order")
	}

	orders, ok := s.Get("orders").(map[string]bool)
	if !ok {
		return errors.New("order not found")
	}

	orders[order.ID] = true

	order.Currency = cart.Currency
	err = initCheckoutItems(&cart, order)
	if err != nil {
		logger.Log(c, err)
		return errors.New("error while adding items to order")
	}

	now := time.Now().Unix()
	order.DateCreated = now
	order.DateUpdated = now
	order.Status = "created"

	err = databases.UpdateOrder(order)
	if err != nil {
		logger.Log(c, err)
		return errors.New("error while updating order")
	}

	s.Set("order_id", order.ID)

	jsonSuccess(c, map[string]any{"order": order})

	return nil
}

func apiGetActiveOrder(c *gin.Context, s sessions.Session) error {
	orderID, ok := s.Get("order_id").(string)
	if !ok {
		err := errors.New("no active order")
		logger.Log(c, err)
		return err
	}

	order, err := databases.FindOrder(orderID)
	if err != nil {
		logger.Log(c, err)
		return errors.New("error while retrieving order")
	}

	if order.Status == "approved" {
		err := fmt.Errorf("error %s is already approved", orderID)
		logger.Log(c, err)
		return err
	}

	jsonSuccess(c, map[string]any{"order": order})

	return nil
}

func initCheckoutItems(cart *model.Cart, order *model.Order) error {
	log.Println(cart.Items)
	for productID, quantity := range cart.Items {
		p, err := databases.GetProduct(productID)
		if err != nil {
			log.Println(err)
			return errors.New("error during order initialization")
		}

		price, err := databases.GetRetailPrice(productID, order.Currency)
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
	address, ok := s.Get("user_infos").(model.Address)
	if !ok {
		err := errors.New("user infos not found")
		logger.Log(c, err)
		return err
	}

	jsonSuccess(c, map[string]any{"user_infos": address})
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
		return errors.New("error while reading param shipping_address")
	}

	if err := checkAddress(&shippingAddress); err != nil {
		log.Println(err)
		return fmt.Errorf("incomplete shipping adress: %w", err)
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
			return errors.New("error while reading param billing_address")
		}

		if err := checkAddress(&billingAddress); err != nil {
			return fmt.Errorf("incomplete billing adress: %w", err)
		}
	}

	orderID, ok := s.Get("order_id").(string)
	if !ok {
		err := errors.New("error while retrieving order id")
		logger.Log(c, err)
		return err
	}

	order, err := databases.FindOrder(orderID)
	if err != nil {
		logger.Log(c, err)
		return errors.New("error while retrieving order")
	}

	if order.Status == "approved" {
		err := fmt.Errorf("error %s is already approved", orderID)
		logger.Log(c, err)
		return err
	}

	order.ShippingAddress = shippingAddress
	order.SameBillingAddress = sameBillingAddress
	order.BillingAddress = billingAddress

	err = databases.UpdateOrder(order)
	if err != nil {
		logger.Log(c, err)
		return errors.New("error while updating order")
	}

	jsonSuccess(c, map[string]any{"order": order})
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

func apiGetShippingMethods(c *gin.Context, s sessions.Session) error {
	orderID, ok := s.Get("order_id").(string)
	if !ok {
		return errors.New("error while retrieving order id")
	}

	order, err := databases.FindOrder(orderID)
	if err != nil {
		logger.Log(c, err)
		return errors.New("error while retrieving order")
	}

	if order.Status == "approved" {
		err := fmt.Errorf("error %s is already approved", orderID)
		logger.Log(c, err)
		return err
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
		p, err := databases.GetProduct(orderItem.ProductID)
		if err != nil {
			logger.Log(c, err)
			return errors.New("error while computing shipping rates")
		}

		variantID, err := strconv.Atoi(p.ExternalID1)
		if err != nil {
			logger.Log(c, err)
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
		logger.Log(c, err)
		return errors.New("error while computing shipping rates")
	}

	log.Println(order)
	order.ShippingInfos = shippingInfos
	for _, shippingInfo := range order.ShippingInfos {
		order.ShippingMethod = shippingInfo.Shipping
		break
	}

	err = computeTaxRate(order)
	if err != nil {
		logger.Log(c, err)
		return errors.New("error while computing tax rate")
	}

	err = databases.UpdateOrder(order)
	if err != nil {
		logger.Log(c, err)
		return errors.New("error while updating order")
	}

	jsonSuccess(c, map[string]any{"order": order})
	return nil
}

func apiSetShippingMethod(c *gin.Context, s sessions.Session, params map[string]any) error {
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

	order, err := databases.FindOrder(orderID)
	if err != nil {
		log.Println(err)
		return errors.New("error while retrieving order")
	}

	if order.Status == "approved" {
		return fmt.Errorf("error %s is already approved", orderID)
	}

	order.ShippingMethod = method
	err = databases.UpdateOrder(order)
	if err != nil {
		log.Println(err)
		return errors.New("error while updating order")
	}

	jsonSuccess(c, map[string]any{"order": order})
	return nil
}

func apiGetOrder(c *gin.Context, s sessions.Session, params map[string]any) error {
	if params == nil {
		return errors.New("no params provided")
	}

	orderID, ok := params["order_id"].(string)
	if !ok {
		return errors.New("invalid order id")
	}

	orders, ok := s.Get("orders").(map[string]bool)
	if !ok {
		orders = make(map[string]bool)
	}

	if !orders[orderID] {
		return errors.New("this order doesn't belong to this user")
	}

	order, err := databases.FindOrder(orderID)
	if err != nil {
		logger.Log(c, err)
		return errors.New("error while getting order")
	}

	jsonSuccess(c, map[string]any{"order": order})
	return nil
}
