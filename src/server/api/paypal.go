package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/plutov/paypal/v4"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/databases/shop"
	"shop.loadout.tf/src/server/email"
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

	order, err := shop.GetOrder(orderID)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	if order.Status == "approved" {
		// Remove the order from the session, as is is already approved
		s.Delete("order_id")
		logger.Log(c, fmt.Errorf("error %s is already approved", orderID))
		return CreateApiError(UnexpectedError)
	}

	client, err := paypal.NewClient(paypalConfig.ClientID, paypalConfig.ClientSecret, paypal.APIBaseSandBox)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	shippingPrice, err := order.GetShippingPrice()
	if err != nil {
		logger.Log(c, fmt.Errorf("unable to get order %s shipping price: %w", orderID, err))
		return CreateApiError(UnexpectedError)
	}

	taxPrice, err := order.GetTaxPrice()
	if err != nil {
		logger.Log(c, fmt.Errorf("unable to get order %s tax price: %w", orderID, err))
		return CreateApiError(UnexpectedError)
	}

	totalPrice, err := order.GetTotalPrice()
	if err != nil {
		logger.Log(c, fmt.Errorf("unable to get order %s total price: %w", orderID, err))
		return CreateApiError(UnexpectedError)
	}

	paypalOrder, err := client.CreateOrder(
		context.Background(),
		paypal.OrderIntentCapture,
		[]paypal.PurchaseUnitRequest{
			{
				Amount: &paypal.PurchaseUnitAmount{
					Value:    totalPrice.String(),
					Currency: order.Currency,
					Breakdown: &paypal.PurchaseUnitAmountBreakdown{
						ItemTotal: &paypal.Money{
							Currency: order.Currency,
							Value:    order.GetItemsPrice().String(),
						},
						Shipping: &paypal.Money{
							Currency: order.Currency,
							Value:    shippingPrice.String(),
						},
						TaxTotal: &paypal.Money{
							Currency: order.Currency,
							Value:    taxPrice.String(),
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
	err = shop.UpdateOrder(order, shop.UpdateOrderFields{Status: true, PaypalOrderID: true})
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, map[string]any{"paypal_order_id": paypalOrder.ID})
	return nil
}

func processPaypalCaptureError(c *gin.Context, params map[string]any, e any) {
	var ok bool
	var content = map[string]any{
		"error":   fmt.Sprint(e),
		"request": c.GetHeader("X-Request-ID"),
	}
	paypalOrderId, ok := params["paypal_order_id"].(string)
	if ok {
		content["paypal_order_id"] = paypalOrderId
	}
	id, err := shop.InsertPaypalError(content)
	if err != nil {
		logger.Log(c, err)
		return
	}
	err = email.SendMail(email.GetMailOrigin(), email.GetMailDestination(), "shop.loadout.tf: paypal capture error "+strconv.FormatInt(id, 10), "")
	if err != nil {
		logger.Log(c, err)
		return
	}
}

func apiCapturePaypalOrder(c *gin.Context, s sessions.Session, params map[string]any) (paypalErr apiError) {
	defer func() {
		if err := recover(); err != nil {
			// We panicked: store the error and return an api error for the regular error handling system
			processPaypalCaptureError(c, params, err)
			paypalErr = CreateApiError(UnexpectedError)
			return
		}

		// Intercept a regular error: store the error and proceed normaly
		if paypalErr != nil {
			processPaypalCaptureError(c, params, paypalErr)
		}

	}()

	if params == nil {
		return CreateApiError(NoParamsError)
	}

	var ok bool
	paypalOrderId, ok := params["paypal_order_id"].(string)
	if !ok {
		return CreateApiError(InvalidParamPaypalOrderID)
	}

	if len(paypalOrderId) > 36 {
		return CreateApiError(InvalidParamPaypalOrderID)
	}
	if !IsAlphaNumeric(paypalOrderId) {
		return CreateApiError(InvalidParamPaypalOrderID)
	}

	client, err := paypal.NewClient(paypalConfig.ClientID, paypalConfig.ClientSecret, paypal.APIBaseSandBox)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	paypalOrder, err := client.GetOrder(
		context.Background(),
		paypalOrderId,
	)

	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	if paypalOrder.Status != "APPROVED" {
		logger.Log(c, errors.New("paypal order is not approved"))
		return CreateApiError(UnexpectedError)
	}

	if len(paypalOrder.PurchaseUnits) < 1 {
		logger.Log(c, fmt.Errorf("paypal order %s don't have purchase units", paypalOrderId))
		return CreateApiError(UnexpectedError)
	}

	// Get the order id from the first and only purchase unit
	purchaseUnit := paypalOrder.PurchaseUnits[0]
	order, err := shop.GetOrder(purchaseUnit.CustomID)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	err = approveOrder(order)
	if err != nil {
		logger.Log(c, fmt.Errorf("error while approving order %s", paypalOrderId))
		return CreateApiError(UnexpectedError)
	}

	var userId string
	authSession := sess.GetAuthSession(c)
	if userId, ok = authSession.Get("user_id").(string); !ok {
		logger.Log(c, fmt.Errorf("error while getting user id from session %s", paypalOrderId))
		return CreateApiError(UnexpectedError)
	}

	err = shop.UserAddOrder(userId, order.ID)
	if err != nil {
		logger.Log(c, fmt.Errorf("error while attaching order %s to user %s", paypalOrderId, userId))
		return CreateApiError(UnexpectedError)
	}

	clearCart(c, s)

	// Remove the order from the session
	s.Delete("order_id")

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
		err := shop.ClearUserCart(userID)
		if err != nil {
			logger.Log(c, err)
		}
	}

}
