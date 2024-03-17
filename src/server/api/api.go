package api

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	appSessions "shop.loadout.tf/src/server/sessions"
)

var _ = registerToken()

func registerToken() bool {
	gob.Register(map[string]interface{}{})
	gob.Register(struct{}{})
	return true
}

type ApiHandler struct {
}

func (handler ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body = make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		jsonError(w, r, errors.New("Bad request"))
		return
	}

	action, ok := body["action"]
	if !ok {
		jsonError(w, r, errors.New("Bad request: no action parameter"))
		return
	}

	session := initSession(w, r)

	params, ok := body["params"]
	var m map[string]interface{}
	if ok {
		m = params.(map[string]interface{})
	}

	switch action {
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

	default:
		jsonError(w, r, NotFoundError{})
		return
	}

	if err != nil {
		jsonError(w, r, err)
	}

	saveSession(w, r, session)
}

func initSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
	session := appSessions.GetSession(r)

	values := session.Values

	if _, ok := values["currency"]; !ok {
		log.Println("Setting currency")
		values["currency"] = "USD" //TODO: set depending on ip
	}

	if _, ok := values["favorites"]; !ok {
		log.Println("Setting favorites")
		values["favorites"] = make(map[string]interface{})
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
		}*/
	return session
}

func saveSession(w http.ResponseWriter, r *http.Request, s *sessions.Session) error {
	err := s.Save(r, w)
	if err != nil {
		log.Println("Error while saving session: ", err)
	}
	return nil
}
