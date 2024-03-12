package api

import (
	"encoding/json"
	"errors"
	_ "log"
	"net/http"
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

	switch action {
	case "get-countries":
		err = getCountries(w, r)
	case "get-currency":
		err = getCurrency(w, r)
	default:
		jsonError(w, r, NotFoundError{})
		return
	}

	if err != nil {
		jsonError(w, r, err)
	}
}
