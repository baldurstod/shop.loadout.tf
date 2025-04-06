package api

import (
	"encoding/gob"
	"errors"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
		log.Println(err)
		jsonError(c, errors.New("bad request"))
		return
	}

	session := initSession(c)
	defer sess.SaveSession(session)

	log.Println("Action: " + request.Action)

	switch request.Action {
	case "get-cart":
		err = getCart(c, session)
	case "get-countries":
		err = getCountries(c)
	case "get-currency":
		err = getCurrency(c, session)
	case "get-favorites":
		err = getFavorites(c, session)
	case "get-product":
		err = getProduct(c, session, request.Params)
	case "get-products":
		err = getProducts(c, session)
	case "get-order":
		err = apiGetOrder(c, session, request.Params)
	case "send-message":
		err = sendMessage(c, request.Params)
	case "set-favorite":
		err = setFavorite(c, session, request.Params)
	case "add-product":
		err = addProduct(c, session, request.Params)
	case "set-product-quantity":
		err = setProductQuantity(c, session, request.Params)
	case "init-checkout":
		err = initCheckout(c, session)
	case "get-active-order":
		err = apiGetActiveOrder(c, session)
	case "create-product":
		err = apiCreateProduct(c, request.Params)
	case "get-user-info":
		err = apiGetUserInfo(c, session)
	case "set-shipping-address":
		err = apiSetShippingAddress(c, session, request.Params)
	case "get-shipping-methods":
		err = apiGetShippingMethods(c, session)
	case "set-shipping-method":
		err = apiSetShippingMethod(c, session, request.Params)
	case "create-paypal-order":
		err = apiCreatePaypalOrder(c, session)
	case "capture-paypal-order":
		err = apiCapturePaypalOrder(c, session, request.Params)
	case "get-printful-products":
		err = apiGetPrintfulProducts(c, request.Params)
	case "get-printful-product":
		err = apiGetPrintfulProduct(c, request.Params)
	case "get-printful-categories":
		err = apiGetPrintfulCategories(c)
	case "get-printful-mockup-styles":
		err = apiGetPrintfulMockupStyles(c, request.Params)
	case "get-printful-product-prices":
		err = apiGetPrintfulProductPrices(c, request.Params)
	case "get-printful-mockup-templates":
		err = apiGetPrintfulMockupTemplates(c, request.Params)
	default:
		jsonError(c, NotFoundError{})
		return
	}

	if err != nil {
		jsonError(c, err)
	}
}

func initSession(c *gin.Context) sessions.Session {
	session := sess.GetSession(c)

	if v := session.Get("currency"); v == nil {
		session.Set("currency", "USD") //TODO: set depending on ip
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

	//session.Values["currency"]
	/*
		req.session.paypal = req.session.paypal ?? {};
		req.session.paypal.token = req.session.paypal.token ?? {};
		req.session.currency = req.session.currency ?? DEFAULT_CURRENCY;//TODO: set depending on ip
		req.session.products = req.session.products ?? {};
		req.session.products.favorites = req.session.products.favorites ?? [];
		req.session.products.visited = req.session.products.visited ?? {};
		req.session.orders = req.session.orders ?? [];
		req.session.communications = req.session.communications ?? [];

		let user;
		if (req.session.userId) {
			user = await this.#restoreUser(req.session.userId);
			req.user = user;
		}

		if (user) {
			req.cart = user.cart;
		} else {
			const cart = new Cart();
			cart.fromJSON(req.session.cart);
			req.cart = cart;
		}*/
	return session
}
