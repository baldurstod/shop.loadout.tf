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
	"shop.loadout.tf/src/server/databases"
	"shop.loadout.tf/src/server/logger"
	"shop.loadout.tf/src/server/model"
)

var paypalConfig config.Paypal

func SetPaypalConfig(config config.Paypal) {
	paypalConfig = config
}

func apiCreatePaypalOrder(c *gin.Context, s sessions.Session) error {
	orderID, ok := s.Get("order_id").(string)
	if !ok {
		err := errors.New("error while retrieving order id")
		logger.Log(c, err)
		return err
	}

	order, err := databases.FindOrder(orderID)
	if err != nil {
		logger.Log(c, err)
		return errors.New("error while retrieving order")
	}

	if order.Status == "approved" {
		err := fmt.Errorf("error %s is already approved", orderID)
		logger.Log(c, err)
		return err
	}

	client, err := paypal.NewClient(paypalConfig.ClientID, paypalConfig.ClientSecret, paypal.APIBaseSandBox)
	if err != nil {
		logger.Log(c, err)
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
		logger.Log(c, err)
		return errors.New("error while creating paypal order")
	}

	log.Println("Got paypal order:", paypalOrder)

	order.PaypalOrderID = paypalOrder.ID
	err = databases.UpdateOrder(order)
	if err != nil {
		logger.Log(c, err)
		return errors.New("error while updating order")
	}

	jsonSuccess(c, map[string]any{"paypal_order_id": paypalOrder.ID})
	return nil
}

func apiCapturePaypalOrder(c *gin.Context, s sessions.Session, params map[string]any) error {
	if params == nil {
		return errors.New("no params provided")
	}

	var id any
	var ok bool
	if id, ok = params["paypal_order_id"]; !ok {
		return errors.New("missing param paypal_order_id")
	}

	orderId, ok := id.(string)
	if !ok {
		return errors.New("param paypal_order_id is not a string")
	}

	if len(orderId) > 36 {
		return errors.New("paypal order id is too long")
	}
	if !IsAlphaNumeric(orderId) {
		return errors.New("paypal order id has a wrong format " + orderId)
	}

	client, err := paypal.NewClient(paypalConfig.ClientID, paypalConfig.ClientSecret, paypal.APIBaseSandBox)
	if err != nil {
		logger.Log(c, err)
		return errors.New("error while creating paypal client")
	}

	paypalOrder, err := client.GetOrder(
		context.Background(),
		orderId,
	)

	if err != nil {
		logger.Log(c, err)
		return errors.New("error while retrieving paypal order")
	}

	if paypalOrder.Status != "APPROVED" {
		err := errors.New("paypal order is not approved")
		logger.Log(c, err)
		return err
	}

	order, err := databases.FindOrderByPaypalID(orderId)
	if err != nil {
		logger.Log(c, err)
		return errors.New("error while retrieving order")
	}

	err = approveOrder(order)
	if err != nil {
		logger.Log(c, err)
		return fmt.Errorf("error while approving order %s", orderId)
	}

	if cart, ok := s.Get("cart").(model.Cart); ok {
		cart.Clear()
	}

	jsonSuccess(c, map[string]any{"order": order})
	return nil
}
