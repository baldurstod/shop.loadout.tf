package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/plutov/paypal/v4"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/databases"
	"shop.loadout.tf/src/server/logger"
	"shop.loadout.tf/src/server/model"
	sess "shop.loadout.tf/src/server/session"
)

var paypalConfig config.Paypal

func SetPaypalConfig(config config.Paypal) {
	paypalConfig = config
}

func apiCreatePaypalOrder(c *gin.Context, s sessions.Session) apiError {
	orderID, ok := s.Get("order_id").(string)
	if !ok {
		logger.Log(c, errors.New("error while retrieving order id"))
		return CreateApiError(UnexpectedError)
	}

	order, err := databases.FindOrder(orderID)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	if order.Status == "approved" {
		logger.Log(c, fmt.Errorf("error %s is already approved", orderID))
		return CreateApiError(UnexpectedError)
	}

	client, err := paypal.NewClient(paypalConfig.ClientID, paypalConfig.ClientSecret, paypal.APIBaseSandBox)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
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
		return CreateApiError(UnexpectedError)
	}

	order.PaypalOrderID = paypalOrder.ID
	err = databases.UpdateOrder(order)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, map[string]any{"paypal_order_id": paypalOrder.ID})
	return nil
}

func apiCapturePaypalOrder(c *gin.Context, s sessions.Session, params map[string]any) apiError {
	if params == nil {
		return CreateApiError(NoParamsError)
	}

	var ok bool
	orderId, ok := params["paypal_order_id"].(string)
	if !ok {
		return CreateApiError(InvalidParamPaypalOrderID)
	}

	if len(orderId) > 36 {
		return CreateApiError(InvalidParamPaypalOrderID)
	}
	if !IsAlphaNumeric(orderId) {
		return CreateApiError(InvalidParamPaypalOrderID)
	}

	client, err := paypal.NewClient(paypalConfig.ClientID, paypalConfig.ClientSecret, paypal.APIBaseSandBox)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	paypalOrder, err := client.GetOrder(
		context.Background(),
		orderId,
	)

	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	if paypalOrder.Status != "APPROVED" {
		logger.Log(c, errors.New("paypal order is not approved"))
		return CreateApiError(UnexpectedError)
	}

	order, err := databases.FindOrderByPaypalID(orderId)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	err = approveOrder(order)
	if err != nil {
		logger.Log(c, fmt.Errorf("error while approving order %s", orderId))
		return CreateApiError(UnexpectedError)
	}

	clearCart(c, s)

	jsonSuccess(c, map[string]any{"order": order})
	return nil
}

func clearCart(c *gin.Context, s sessions.Session) {
	// Clear cart in session
	if cart, ok := s.Get("cart").(model.Cart); ok {
		cart.Clear()
		s.Set("cart", cart)
		s.Save()
	}

	// Clear user cart
	authSession := sess.GetAuthSession(c)
	if userID, ok := authSession.Get("user_id").(string); ok {
		err := databases.ClearUserCart(userID)
		if err != nil {
			logger.Log(c, err)
		}
	}

}
