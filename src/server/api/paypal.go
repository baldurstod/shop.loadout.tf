package api

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/plutov/paypal/v4"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/model"
	"shop.loadout.tf/src/server/mongo"
)

var paypalConfig config.Paypal

func SetPaypalConfig(config config.Paypal) {
	paypalConfig = config
}

func apiCreatePaypalOrder(c *gin.Context, s sessions.Session) error {
	//log.Println(s)

	orderID, ok := s.Get("order_id").(string)
	if !ok {
		return errors.New("error while retrieving order id")
	}

	order, err := mongo.FindOrder(orderID)
	if err != nil {
		log.Println(err)
		return errors.New("error while retrieving order")
	}

	if order.Status == "approved" {
		return fmt.Errorf("error %s is already approved", orderID)
	}

	fmt.Println(order)

	client, err := paypal.NewClient(paypalConfig.ClientID, paypalConfig.ClientSecret, paypal.APIBaseSandBox)
	if err != nil {
		log.Println(err)
		return errors.New("error while creating paypal client")
	}

	paypalOrder, err := client.CreateOrder(
		context.Background(),
		paypal.OrderIntentCapture,
		[]paypal.PurchaseUnitRequest{
			{
				Amount: &paypal.PurchaseUnitAmount{
					Value:    order.GetTotalPrice().String(),
					Currency: order.Currency,
					Breakdown: &paypal.PurchaseUnitAmountBreakdown{
						ItemTotal: &paypal.Money{
							Currency: order.Currency,
							Value:    order.GetItemsPrice().String(),
						},
						Shipping: &paypal.Money{
							Currency: order.Currency,
							Value:    order.GetShippingPrice().String(),
						},
						TaxTotal: &paypal.Money{
							Currency: order.Currency,
							Value:    order.GetTaxPrice().String(),
						},
					},
					/*
						amount: {
							currency_code: currency,
							value: roundPrice(currency, order.totalPrice),
							breakdown: {
							}
						},
					*/
				},
				CustomID: order.ID,
				Shipping: &paypal.ShippingDetail{
					Name: &paypal.Name{
						FullName: order.ShippingAddress.GetFullName(),
					},
					Address: &paypal.ShippingDetailAddressPortable{
						AddressLine1: order.ShippingAddress.Address1,
						AddressLine2: order.ShippingAddress.Address2,
						AdminArea1:   order.ShippingAddress.StateCode,
						AdminArea2:   order.ShippingAddress.City,
						PostalCode:   order.ShippingAddress.PostalCode,
						CountryCode:  order.ShippingAddress.CountryCode,
					},
				},
			},
		},
		&paypal.CreateOrderPayer{},
		&paypal.ApplicationContext{
			ShippingPreference: paypal.ShippingPreferenceSetProvidedAddress,
		},
	)

	if err != nil {
		log.Println(err)
		return errors.New("error while creating paypal order")
	}

	log.Println("Got paypal order:", paypalOrder)

	order.PaypalOrderID = paypalOrder.ID
	err = mongo.UpdateOrder(order)
	if err != nil {
		log.Println(err)
		return errors.New("error while updating order")
	}

	jsonSuccess(c, map[string]interface{}{"paypal_order_id": paypalOrder.ID})
	return nil
}

func apiCapturePaypalOrder(c *gin.Context, s sessions.Session, params map[string]interface{}) error {
	if params == nil {
		return errors.New("no params provided")
	}

	var id interface{}
	var ok bool
	if id, ok = params["paypal_order_id"]; !ok {
		return errors.New("missing param paypal_order_id")
	}

	orderId := id.(string)

	if len(orderId) > 36 {
		return errors.New("paypal order id is too long")
	}
	if !IsAlphaNumeric(orderId) {
		return errors.New("paypal order id has a wrong format " + orderId)
	}

	client, err := paypal.NewClient(paypalConfig.ClientID, paypalConfig.ClientSecret, paypal.APIBaseSandBox)
	if err != nil {
		log.Println(err)
		return errors.New("error while creating paypal client")
	}

	paypalOrder, err := client.GetOrder(
		context.Background(),
		orderId,
	)

	if err != nil {
		log.Println(err)
		return errors.New("error while retrieving paypal order")
	}

	if paypalOrder.Status != "APPROVED" {
		return errors.New("paypal order is not approved")
	}

	order, err := mongo.FindOrderByPaypalID(orderId)
	if err != nil {
		log.Println(err)
		return errors.New("error while retrieving order")
	}

	err = approveOrder(order)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error while approving order %s", orderId)
	}

	cart := s.Get("cart").(model.Cart)
	cart.Clear()

	jsonSuccess(c, map[string]interface{}{"order": order})
	return nil
}
