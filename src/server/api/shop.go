package api

import (
	"errors"
	"fmt"
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
	sess "shop.loadout.tf/src/server/session"
)

func apiGetCurrency(c *gin.Context, s sessions.Session) apiError {
	currency, ok := s.Get("currency").(string)
	if !ok {
		currency = constants.DEFAULT_CURRENCY
	}

	jsonSuccess(c, map[string]any{"currency": currency})
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

func apiGetProduct(c *gin.Context, s sessions.Session, params map[string]any) apiError {
	if params == nil {
		return CreateApiError(NoParamsError)
	}

	currency, ok := s.Get("currency").(string)
	if !ok {
		currency = constants.DEFAULT_CURRENCY
	}
	prices := NewProductPrice(currency)

	productID, ok := params["product_id"].(string)
	if !ok {
		return CreateApiError(InvalidParamProductID)
	}

	product, err := databases.FindProduct(productID)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
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
			continue
		}
		prices.Prices[id] = price.RetailPrice.String()
	}

	jsonSuccess(c, map[string]any{
		"product": product,
		"prices":  prices,
	})
	return nil
}

func apiGetProducts(c *gin.Context, s sessions.Session) apiError {
	p, err := databases.GetProducts()
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

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
			continue
		}
		prices.Prices[id] = price.RetailPrice.String()
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

func apiSendMessage(c *gin.Context, params map[string]any) apiError {
	if params == nil {
		return CreateApiError(NoParamsError)
	}

	subject, ok := params["subject"].(string)
	if !ok || subject == "" || len(subject) < 5 {
		return CreateApiError(InvalidParamSubject)
	}

	email, ok := params["email"].(string)
	if !ok || email == "" || !validEmail(email) {
		return CreateApiError(InvalidParamEmail)
	}

	content, ok := params["content"].(string)
	if !ok || content == "" || len(content) < 10 {
		return CreateApiError(InvalidParamContent)
	}

	id, err := databases.SendContact(subject, email, content)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, map[string]any{"message_id": id})
	return nil
}

func apiAddProduct(c *gin.Context, s sessions.Session, params map[string]any) apiError {
	if params == nil {
		return CreateApiError(NoParamsError)
	}

	productId, ok := params["product_id"].(string)
	if !ok {
		return CreateApiError(InvalidParamProductID)
	}

	quantity, ok := params["quantity"].(float64)
	if !ok {
		return CreateApiError(InvalidParamQuantity)
	}

	cart, ok := s.Get("cart").(model.Cart)
	if !ok {
		logger.Log(c, errors.New("cart not found"))
		return CreateApiError(UnexpectedError)
	}

	cart.AddQuantity(productId, uint(quantity))
	s.Delete("order_id")

	authSession := sess.GetAuthSession(c)
	if userID, ok := authSession.Get("user_id").(string); ok {
		user, err := databases.FindUserByID(userID)
		if err != nil {
			logger.Log(c, err)
		} else {
			cart = user.Cart
			cart.AddQuantity(productId, uint(quantity))
			err = databases.SetUserCart(userID, cart)
			if err != nil {
				logger.Log(c, err)
			}
		}
	}

	jsonSuccess(c, map[string]any{"cart": cart})
	return nil
}

func apiSetProductQuantity(c *gin.Context, s sessions.Session, params map[string]any) apiError {
	if params == nil {
		return CreateApiError(NoParamsError)
	}

	productId, ok := params["product_id"].(string)
	if !ok {
		return CreateApiError(InvalidParamProductID)
	}

	quantity, ok := params["quantity"].(float64)
	if !ok {
		return CreateApiError(InvalidParamQuantity)
	}

	cart, ok := s.Get("cart").(model.Cart)
	if !ok {
		logger.Log(c, errors.New("cart not found"))
		return CreateApiError(UnexpectedError)
	}

	cart.SetQuantity(productId, uint(quantity))
	s.Delete("order_id")

	authSession := sess.GetAuthSession(c)
	if userID, ok := authSession.Get("user_id").(string); ok {
		user, err := databases.FindUserByID(userID)
		if err != nil {
			logger.Log(c, err)
		} else {
			cart = user.Cart
			cart.SetQuantity(productId, uint(quantity))
			err = databases.SetUserCart(userID, cart)
			if err != nil {
				logger.Log(c, err)
			}
		}
	}

	jsonSuccess(c, map[string]any{"cart": cart})
	return nil
}

func apiGetCart(c *gin.Context, s sessions.Session) apiError {
	authSession := sess.GetAuthSession(c)
	if userID, ok := authSession.Get("user_id").(string); ok {
		user, err := databases.FindUserByID(userID)
		if err != nil {
			logger.Log(c, err)
			return CreateApiError(UnexpectedError)
		}
		jsonSuccess(c, map[string]any{"cart": user.Cart})
		return nil
	}

	cart, ok := s.Get("cart").(model.Cart)
	if !ok {
		cart = model.NewCart()
	}

	jsonSuccess(c, map[string]any{"cart": cart})
	return nil
}

func apiInitCheckout(c *gin.Context, s sessions.Session) apiError {
	cart, ok := s.Get("cart").(model.Cart)
	if !ok {
		logger.Log(c, errors.New("cart not found"))
		return CreateApiError(UnexpectedError)
	}

	order, err := databases.CreateOrder()
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	orders, ok := s.Get("orders").(map[string]bool)
	if !ok {
		return CreateApiError(UnexpectedError)
	}

	orders[order.ID] = true

	order.Currency = cart.Currency
	err = initCheckoutItems(&cart, order)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	now := time.Now().Unix()
	order.DateCreated = now
	order.DateUpdated = now
	order.Status = "created"

	err = databases.UpdateOrder(order)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	s.Set("order_id", order.ID)

	jsonSuccess(c, map[string]any{"order": order})

	return nil
}

func apiGetActiveOrder(c *gin.Context, s sessions.Session) apiError {
	orderID, ok := s.Get("order_id").(string)
	if !ok {
		logger.Log(c, errors.New("no active order"))
		return CreateApiError(UnexpectedError)
	}

	order, err := databases.FindOrder(orderID)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	if order.Status == "approved" {
		logger.Log(c, fmt.Errorf("error %s is already approved", orderID))
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, map[string]any{"order": order})

	return nil
}

func initCheckoutItems(cart *model.Cart, order *model.Order) error {
	for productID, quantity := range cart.Items {
		p, err := databases.GetProduct(productID)
		if err != nil {
			return fmt.Errorf("error while getting product %s: %w", productID, err)
		}

		price, err := databases.GetRetailPrice(productID, order.Currency)
		if err != nil {
			return fmt.Errorf("error while getting retail price for product %s: %w", productID, err)
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

	return nil
}

func apiGetUserInfo(c *gin.Context, s sessions.Session) apiError {
	address, ok := s.Get("user_infos").(model.Address)
	if !ok {
		logger.Log(c, errors.New("user infos not found"))
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, map[string]any{"user_infos": address})
	return nil
}

func apiSetShippingAddress(c *gin.Context, s sessions.Session, params map[string]any) apiError {
	shippingAddress := model.Address{}
	billingAddress := model.Address{}

	if params["shipping_address"] == nil {
		return CreateApiError(InvalidParamShippingAddress)
	}

	err := mapstructure.Decode(params["shipping_address"], &shippingAddress)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(InvalidParamShippingAddress)
	}

	if err := checkAddress(&shippingAddress); err != nil {
		logger.Log(c, fmt.Errorf("incomplete shipping adress: %w", err))
		return CreateApiError(InvalidParamShippingAddress)
	}

	sameBillingAddress, ok := params["same_billing_address"].(bool)
	if !ok {
		return CreateApiError(InvalidParamSameBillingAddress)
	}

	if !sameBillingAddress {
		if params["billing_address"] == nil {
			return CreateApiError(InvalidParamBillingAddress)
		}

		err := mapstructure.Decode(params["billing_address"], &billingAddress)
		if err != nil {
			logger.Log(c, err)
			return CreateApiError(InvalidParamBillingAddress)
		}

		if err := checkAddress(&billingAddress); err != nil {
			logger.Log(c, fmt.Errorf("incomplete billing adress: %w", err))
			return CreateApiError(InvalidParamBillingAddress)
		}
	}

	orderID, ok := s.Get("order_id").(string)
	if !ok {
		logger.Log(c, errors.New("error while retrieving order id"))
		return CreateApiError(UnexpectedError)
	}

	order, err := databases.FindOrder(orderID)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	if order.Status == "approved" {
		logger.Log(c, fmt.Errorf("error %s is already approved", orderID))
		return CreateApiError(UnexpectedError)
	}

	order.ShippingAddress = shippingAddress
	order.SameBillingAddress = sameBillingAddress
	order.BillingAddress = billingAddress

	err = databases.UpdateOrder(order)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
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

func apiGetShippingMethods(c *gin.Context, s sessions.Session) apiError {
	orderID, ok := s.Get("order_id").(string)
	if !ok {
		logger.Log(c, errors.New("error while retrieving order id"))
		return CreateApiError(UnexpectedError)
	}

	order, err := databases.FindOrder(orderID)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	if order.Status == "approved" {
		logger.Log(c, fmt.Errorf("error %s is already approved", orderID))
		return CreateApiError(UnexpectedError)
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
			return CreateApiError(UnexpectedError)
		}

		variantID, err := strconv.Atoi(p.ExternalID1)
		if err != nil {
			logger.Log(c, fmt.Errorf("unable to convert external id %s: %w", p.ExternalID1, err))
			return CreateApiError(UnexpectedError)
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
		return CreateApiError(UnexpectedError)
	}

	order.ShippingInfos = shippingInfos
	for _, shippingInfo := range order.ShippingInfos {
		order.ShippingMethod = shippingInfo.Shipping
		break
	}

	err = computeTaxRate(order)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	err = databases.UpdateOrder(order)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, map[string]any{"order": order})
	return nil
}

func apiSetShippingMethod(c *gin.Context, s sessions.Session, params map[string]any) apiError {
	method, ok := params["method"].(string)
	if !ok {
		return CreateApiError(InvalidParamMethod)
	}

	orderID, ok := s.Get("order_id").(string)
	if !ok {
		logger.Log(c, errors.New("error while retrieving order id"))
		return CreateApiError(UnexpectedError)
	}

	order, err := databases.FindOrder(orderID)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	if order.Status == "approved" {
		logger.Log(c, fmt.Errorf("error %s is already approved", orderID))
		return CreateApiError(UnexpectedError)
	}

	order.ShippingMethod = method
	err = databases.UpdateOrder(order)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, map[string]any{"order": order})
	return nil
}

func apiGetOrder(c *gin.Context, s sessions.Session, params map[string]any) apiError {
	if params == nil {
		return CreateApiError(NoParamsError)
	}

	orderID, ok := params["order_id"].(string)
	if !ok {
		return CreateApiError(InvalidParamOrderID)
	}

	orders, ok := s.Get("orders").(map[string]bool)
	if !ok {
		orders = make(map[string]bool)
	}

	if !orders[orderID] {
		logger.Log(c, fmt.Errorf("order %s doesn't belong to user with session %s", orderID, s.ID()))
		return CreateApiError(UnexpectedError)
	}

	order, err := databases.FindOrder(orderID)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, map[string]any{"order": order})
	return nil
}
