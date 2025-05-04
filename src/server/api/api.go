package api

import (
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"shop.loadout.tf/src/server/constants"
	"shop.loadout.tf/src/server/logger"
	"shop.loadout.tf/src/server/model"
	sess "shop.loadout.tf/src/server/session"
)

var _ = registerToken()

func registerToken() bool {
	gob.Register(map[string]any{})
	gob.Register(map[string]bool{})
	gob.Register(struct{}{})
	gob.Register(model.Cart{})
	gob.Register(model.Address{})
	return true
}

/*
type ApiHandler struct {
}*/

type ApiRequest struct {
	Action  string         `json:"action" binding:"required"`
	Version int            `json:"version" binding:"required"`
	Params  map[string]any `json:"params"`
}

func ApiHandler(c *gin.Context) {
	var request ApiRequest
	var err error

	if err = c.ShouldBindJSON(&request); err != nil {
		logger.Log(c, err)
		jsonError(c, errors.New("bad request"))
		return
	}

	session := initSession(c)
	defer func() {
		if err := session.Save(); err != nil {
			logger.Log(c, fmt.Errorf("error while saving session: %w", err))
		}
	}()

	var apiError apiError
	switch request.Action {
	case "get-cart":
		apiError = apiGetCart(c, session)
	case "get-countries":
		apiError = apiGetCountries(c)
	case "get-currency":
		apiError = apiGetCurrency(c, session)
	case "get-favorites":
		apiError = apiGetFavorites(c, session)
	case "get-product":
		apiError = apiGetProduct(c, session, request.Params)
	case "get-products":
		apiError = apiGetProducts(c, session)
	case "get-order":
		apiError = apiGetOrder(c, session, request.Params)
	case "send-message":
		apiError = apiSendMessage(c, request.Params)
	case "set-favorite":
		apiError = apiSetFavorite(c, session, request.Params)
	case "add-product":
		apiError = apiAddProduct(c, session, request.Params)
	case "set-product-quantity":
		apiError = apiSetProductQuantity(c, session, request.Params)
	case "init-checkout":
		apiError = apiInitCheckout(c, session)
	case "get-active-order":
		apiError = apiGetActiveOrder(c, session)
	case "create-product":
		apiError = apiCreateProduct(c, request.Params)
	case "get-user-info":
		apiError = apiGetUserInfo(c, session)
	case "set-shipping-address":
		apiError = apiSetShippingAddress(c, session, request.Params)
	case "get-shipping-methods":
		apiError = apiGetShippingMethods(c, session)
	case "set-shipping-method":
		apiError = apiSetShippingMethod(c, session, request.Params)
	case "create-paypal-order":
		apiError = apiCreatePaypalOrder(c, session)
	case "capture-paypal-order":
		apiError = apiCapturePaypalOrder(c, session, request.Params)
	case "create-account":
		apiError = apiCreateAccount(c, session, request.Params)
	case "login":
		apiError = apiLogin(c, session, request.Params)
	case "logout":
		apiError = apiLogout(c, session)
	case "get-printful-products":
		apiError = apiGetPrintfulProducts(c, request.Params)
	case "get-printful-product":
		apiError = apiGetPrintfulProduct(c, request.Params)
	case "get-printful-categories":
		apiError = apiGetPrintfulCategories(c)
	case "get-printful-mockup-styles":
		apiError = apiGetPrintfulMockupStyles(c, request.Params)
	case "get-printful-product-prices":
		apiError = apiGetPrintfulProductPrices(c, request.Params)
	case "get-printful-mockup-templates":
		apiError = apiGetPrintfulMockupTemplates(c, request.Params)
	default:
		jsonError(c, NotFoundError{})
		return
	}

	if apiError != nil {
		jsonError(c, apiError)
	}
}

func initSession(c *gin.Context) sessions.Session {
	session := sess.GetRegularSession(c)

	if v := session.Get("currency"); v == nil {
		session.Set("currency", constants.DEFAULT_CURRENCY) //TODO: set depending on ip
	}

	if v := session.Get("favorites"); v == nil {
		session.Set("favorites", make(map[string]any))
	}

	if v := session.Get("cart"); v == nil {
		session.Set("cart", model.NewCart())
	}

	if v := session.Get("user_infos"); v == nil {
		session.Set("user_infos", model.Address{})
	}

	if v := session.Get("orders"); v == nil {
		session.Set("orders", make(map[string]bool))
	}

	return session
}
