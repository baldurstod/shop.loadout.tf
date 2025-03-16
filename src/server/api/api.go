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
	gob.Register(map[string]interface{}{})
	gob.Register(struct{}{})
	gob.Register(model.Cart{})
	gob.Register(model.Address{})
	return true
}

/*
type ApiHandler struct {
}*/

type ApiRequest struct {
	Action  string                 `json:"action" binding:"required"`
	Version int                    `json:"version" binding:"required"`
	Params  map[string]interface{} `json:"params"`
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
		err = getProduct(c, request.Params)
	case "get-products":
		err = getProducts(c)
	case "get-order":
		err = apiGetOrder(c, request.Params)
	case "send-contact":
		err = sendContact(c, request.Params)
	case "set-favorite":
		err = setFavorite(c, session, request.Params)
	case "add-product":
		err = addProduct(c, session, request.Params)
	case "set-product-quantity":
		err = setProductQuantity(c, session, request.Params)
	case "init-checkout":
		err = initCheckout(c, session)
	case "create-product":
		err = apiCreateProduct(c, request.Params)
	case "get-user-info":
		err = apiGetUserInfo(c, session)
	case "set-shipping-address":
		err = apiSetShippingAddress(c, session, request.Params)
	case "set-shipping-method":
		err = apiSetShippingMethod(c, session, request.Params)
	case "create-paypal-order":
		err = apiCreatePaypalOrder(c, session)
	case "capture-paypal-order":
		err = apiCapturePaypalOrder(c, session, request.Params)
	case "get-printful-products":
		err = apiGetPrintfulProducts(c)
	case "get-printful-product":
		err = apiGetPrintfulProduct(c, request.Params)

	default:
		jsonError(c, NotFoundError{})
		return
	}

	if err != nil {
		jsonError(c, err)
	}
}

/*
	func ServeHTTP(w http.ResponseWriter, r *http.Request) {
		var body = make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			//jsonError(w, r, errors.New("Bad request"))
			return
		}

		action, ok := body["action"]
		if !ok {
			//jsonError(w, r, errors.New("Bad request: no action parameter"))
			return
		}

		session := initSession(w, r)

		params, ok := body["params"]
		var m map[string]interface{}
		if ok {
			m = params.(map[string]interface{})
		}

		switch action {
		case "get-cart":
			err = getCart(w, r, session)
		case "get-countries":
			err = getCountries(w, r)
		case "get-currency":
			err = getCurrency(w, r, session)
		case "get-favorites":
			err = getFavorites(w, r, session)
		case "get-product":
			err = getProduct(w, r, m)
		case "get-products":
			err = getProducts(w, r, session)
		case "send-contact":
			err = sendContact(w, r, m)
		case "set-favorite":
			err = setFavorite(w, r, session, m)
		case "add-product":
			err = addProduct(w, r, session, m)
		case "set-product-quantity":
			err = setProductQuantity(w, r, session, m)
		case "init-checkout":
			err = initCheckout(w, r, session, m)
		case "create-product":
			err = apiCreateProduct(w, r, session, m)
		case "get-user-info":
			err = apiGetUserInfo(w, r, session, m)
		case "set-shipping-address":
			err = apiSetShippingAddress(w, r, session, m)
		case "set-shipping-method":
			err = apiSetShippingMethod(w, r, session, m)
		case "create-paypal-order":
			err = apiCreatePaypalOrder(w, r, session, m)
		default:
			jsonError(w, r, NotFoundError{})
			return
		}

		if err != nil {
			jsonError(w, r, err)
		}
	}*/

func initSession(c *gin.Context) sessions.Session {
	session := sess.GetSession(c)

	//values := session.Values

	if v := session.Get("currency"); v == nil {
		session.Set("currency", "USD") //TODO: set depending on ip
	}

	if v := session.Get("favorites"); v == nil {
		session.Set("favorites", make(map[string]interface{}))
	}

	if v := session.Get("cart"); v == nil {
		session.Set("cart", model.NewCart())
	}

	if v := session.Get("user_infos"); v == nil {
		session.Set("user_infos", model.Address{})
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
