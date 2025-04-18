package api

import "errors"

type NotFoundError struct{}

func (e NotFoundError) Error() string {
	return "Not found"
}

type ApiErrorCode int

const (
	AuthenticationError ApiErrorCode = iota
	NoParamsError
	InvalidParams
	InvalidParamProduct
	InvalidParamProductID
	InvalidParamQuantity
	InvalidParamOrderID
	InvalidParamSubject
	InvalidParamEmail
	InvalidParamIsFavorite
	InvalidParamShippingAddress
	InvalidParamBillingAddress
	InvalidParamSameBillingAddress
	InvalidParamMethod
	InvalidParamPaypalOrderID
	InvalidParamPassword
	InvalidParamCurrency
	InvalidParamContent
	ProductNotFound
	UnexpectedError
)

var apiErrorValues = map[ApiErrorCode]error{
	AuthenticationError:            errors.New("authentication error"),
	NoParamsError:                  errors.New("no params provided"),
	InvalidParams:                  errors.New("invalid parameters"),
	InvalidParamProduct:            errors.New("invalid param product"),
	InvalidParamProductID:          errors.New("invalid param product_id"),
	InvalidParamQuantity:           errors.New("invalid param quantity"),
	InvalidParamOrderID:            errors.New("invalid param order_id"),
	InvalidParamSubject:            errors.New("invalid param subject"),
	InvalidParamEmail:              errors.New("invalid param email"),
	InvalidParamIsFavorite:         errors.New("invalid param is_favorite"),
	InvalidParamShippingAddress:    errors.New("invalid param shipping_address"),
	InvalidParamBillingAddress:     errors.New("invalid param billing_address"),
	InvalidParamSameBillingAddress: errors.New("invalid param same_billing_address"),
	InvalidParamMethod:             errors.New("invalid param method"),
	InvalidParamPaypalOrderID:      errors.New("invalid param paypal_order_id"),
	InvalidParamPassword:           errors.New("invalid param password"),
	InvalidParamCurrency:           errors.New("invalid param currency"),
	InvalidParamContent:            errors.New("invalid param content"),
	ProductNotFound:                errors.New("product not found"),
	UnexpectedError:                errors.New("unexpected error, contact support"),
}

type apiError interface {
	Error() string
	isApiError() bool
}

type apiError2 struct {
	StatusCode int
	Err        error
}

func (e apiError2) Error() string {
	return e.Err.Error()
}

func (e apiError2) isApiError() bool {
	return true
}

func CreateApiError(c ApiErrorCode) apiError2 {
	e := apiErrorValues[c]
	return apiError2{Err: e}
}
