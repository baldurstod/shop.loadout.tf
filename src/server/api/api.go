package api

import (
	"encoding/gob"
	"errors"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"shop.loadout.tf/src/server/model"
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
		jsonError(c, errors.New("Bad request"))
		return
	}

	session := sessions.Default(c)

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
		err = getProducts(c, session)
	case "send-contact":
		err = sendContact(c, request.Params)
	case "set-favorite":
		err = setFavorite(c, session, request.Params)
	case "add-product":
		err = addProduct(c, session, request.Params)
	case "set-product-quantity":
		err = setProductQuantity(c, session, request.Params)
	case "init-checkout":
		err = initCheckout(c, session, request.Params)
	case "create-product":
		err = apiCreateProduct(c, session, request.Params)
	case "get-user-info":
		err = apiGetUserInfo(c, session, request.Params)
	case "set-shipping-address":
		err = apiSetShippingAddress(c, session, request.Params)
	case "set-shipping-method":
		err = apiSetShippingMethod(c, session, request.Params)
	case "create-paypal-order":
		err = apiCreatePaypalOrder(c, session, request.Params)
	default:
		jsonError(c, NotFoundError{})
		return
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
	}

	func initSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
		session := appSessions.GetSession(r)

		values := session.Values

		if _, ok := values["currency"]; !ok {
			values["currency"] = "USD" //TODO: set depending on ip
		}

		if _, ok := values["favorites"]; !ok {
			values["favorites"] = make(map[string]interface{})
		}

		if _, ok := values["cart"]; !ok {
			values["cart"] = model.NewCart()
		}

		if _, ok := values["user_infos"]; !ok {
			values["user_infos"] = model.Address{}
		}

		saveSession(w, r, session)

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
			}* /
		return session
	}
*/
func saveSession(c *gin.Context, s sessions.Session) error {
	err := s.Save()
	if err != nil {
		log.Println("Error while saving session: ", err)
	}
	return nil
}
