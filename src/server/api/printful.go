package api

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/baldurstod/printful-api-model"
	"log"
	"net/http"
	"net/url"
	"shop.loadout.tf/src/server/config"
	//"shop.loadout.tf/src/server/sessions"
	"github.com/gorilla/sessions"
)

var printfulConfig config.Printful

func SetPrintfulConfig(config config.Printful) {
	printfulConfig = config
	log.Println(config)
}

func getCountries(w http.ResponseWriter, r *http.Request) error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	u, err := url.JoinPath(printfulConfig.Endpoint, "/countries")
	if err != nil {
		return errors.New("Error while getting printful url")
	}

	resp, err := http.Get(u)

	if err != nil {
		log.Println(err)
		return errors.New("Error while calling printful api")
	}

	countriesResponse := model.CountriesResponse{}
	err = json.NewDecoder(resp.Body).Decode(&countriesResponse)
	if err != nil {
		return errors.New("Error while decoding printful response")
	}

	jsonSuccess(w, r, countriesResponse.Countries)

	return nil
}

func getCurrency(w http.ResponseWriter, r *http.Request, s *sessions.Session) error {
	//session := sessions.GetSession(r)

	jsonSuccess(w, r, s.Values["currency"])

	return nil
}
