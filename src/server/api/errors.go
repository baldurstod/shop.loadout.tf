package api

import (
	"errors"
	"log"
)

type ApiErrorCode int

const (
	NotFoundError ApiErrorCode = iota
	BadRequestError
	AuthenticationError
	NoParamsError
	InvalidParams
	InvalidParamProduct
	InvalidParamProductID
	InvalidParamQuantity
	InvalidParamOrderID
	InvalidParamSubject
	InvalidParamUsername
	InvalidParamEmail
	InvalidParamCode
	InvalidParamIsFavorite
	InvalidParamShippingAddress
	InvalidParamBillingAddress
	InvalidParamSameBillingAddress
	InvalidParamMethod
	InvalidParamPaypalOrderID
	InvalidParamPassword
	InvalidParamCurrency
	InvalidParamContent
	InvalidParamLanguage
	UnexpectedError
	NotAuthenticated
	AlreadyAuthenticated
	InvalidVerificationCode
	ExpiredVerificationCode
	VerifiedEmailRequired
)

var apiErrorValues = map[ApiErrorCode]error{
	NotFoundError:                  errors.New("not found"),
	BadRequestError:                errors.New("bad request"),
	AuthenticationError:            errors.New("authentication error"),
	NoParamsError:                  errors.New("no params provided"),
	InvalidParams:                  errors.New("invalid parameters"),
	InvalidParamProduct:            errors.New("invalid param product"),
	InvalidParamProductID:          errors.New("invalid param product_id"),
	InvalidParamQuantity:           errors.New("invalid param quantity"),
	InvalidParamOrderID:            errors.New("invalid param order_id"),
	InvalidParamSubject:            errors.New("invalid param subject"),
	InvalidParamUsername:           errors.New("invalid param username"),
	InvalidParamEmail:              errors.New("invalid param email"),
	InvalidParamCode:               errors.New("invalid param code"),
	InvalidParamIsFavorite:         errors.New("invalid param is_favorite"),
	InvalidParamShippingAddress:    errors.New("invalid param shipping_address"),
	InvalidParamBillingAddress:     errors.New("invalid param billing_address"),
	InvalidParamSameBillingAddress: errors.New("invalid param same_billing_address"),
	InvalidParamMethod:             errors.New("invalid param method"),
	InvalidParamPaypalOrderID:      errors.New("invalid param paypal_order_id"),
	InvalidParamPassword:           errors.New("invalid param password"),
	InvalidParamCurrency:           errors.New("invalid param currency"),
	InvalidParamContent:            errors.New("invalid param content"),
	InvalidParamLanguage:           errors.New("invalid param language"),
	UnexpectedError:                errors.New("unexpected error, contact support"),
	NotAuthenticated:               errors.New("user not authenticated"),
	AlreadyAuthenticated:           errors.New("user already authenticated"),
	InvalidVerificationCode:        errors.New("invalid verification code"),
	ExpiredVerificationCode:        errors.New("expired verification code"),
	VerifiedEmailRequired:          errors.New("verified email required"),
}

var apiErrorI18n = map[ApiErrorCode]string{
	NotFoundError:                  "#api_error_not_found",
	BadRequestError:                "#api_error_bad_request",
	AuthenticationError:            "#api_error_authentication_error",
	NoParamsError:                  "#api_error_no_params_provided",
	InvalidParams:                  "#api_error_invalid_parameters",
	InvalidParamProduct:            "#api_error_invalid_param_product",
	InvalidParamProductID:          "#api_error_invalid_param_product_id",
	InvalidParamQuantity:           "#api_error_invalid_param_quantity",
	InvalidParamOrderID:            "#api_error_invalid_param_order_id",
	InvalidParamSubject:            "#api_error_invalid_param_subject",
	InvalidParamUsername:           "#api_error_invalid_param_username",
	InvalidParamEmail:              "#api_error_invalid_param_email",
	InvalidParamCode:               "#api_error_invalid_param_code",
	InvalidParamIsFavorite:         "#api_error_invalid_param_is_favorite",
	InvalidParamShippingAddress:    "#api_error_invalid_param_shipping_address",
	InvalidParamBillingAddress:     "#api_error_invalid_param_billing_address",
	InvalidParamSameBillingAddress: "#api_error_invalid_param_same_billing_address",
	InvalidParamMethod:             "#api_error_invalid_param_method",
	InvalidParamPaypalOrderID:      "#api_error_invalid_param_paypal_order_id",
	InvalidParamPassword:           "#api_error_invalid_param_password",
	InvalidParamCurrency:           "#api_error_invalid_param_currency",
	InvalidParamContent:            "#api_error_invalid_param_content",
	InvalidParamLanguage:           "#api_error_invalid_param_language",
	UnexpectedError:                "#api_error_unexpected_error",
	NotAuthenticated:               "#api_error_user_not_authenticated",
	AlreadyAuthenticated:           "#api_error_user_already_authenticated",
	InvalidVerificationCode:        "#api_error_invalid_verification_code",
	ExpiredVerificationCode:        "#api_error_expired_verification_code",
	VerifiedEmailRequired:          "#api_error_verified_email_required",
}

type apiError interface {
	Error() string
	I18n() string
	isApiError() bool
}

type apiError2 struct {
	StatusCode int
	Err        error
	i18n       string
}

func (e apiError2) Error() string {
	return e.Err.Error()
}

func (e apiError2) I18n() string {
	return e.i18n
}

func (e apiError2) isApiError() bool {
	return true
}

func CreateApiError(c ApiErrorCode) apiError2 {
	e, found := apiErrorValues[c]
	if !found {
		log.Println("Missing message for error code ", c)
		e = apiErrorValues[UnexpectedError]
	}

	i, found := apiErrorI18n[c]
	if !found {
		log.Println("Missing i18n for error code ", c)
		e = apiErrorValues[UnexpectedError]
	}

	return apiError2{Err: e, i18n: i}
}
