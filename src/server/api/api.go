package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"github.com/gorilla/sessions"
	appSessions "shop.loadout.tf/src/server/sessions"
)

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

	switch action {
	case "get-countries":
		err = getCountries(w, r)
	case "get-currency":
		err = getCurrency(w, r, session)
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
		values["currency"] = "USD"//TODO: set depending on ip
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
