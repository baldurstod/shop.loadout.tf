package api

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"shop.loadout.tf/src/server/databases/shop"
	"shop.loadout.tf/src/server/logger"
	"shop.loadout.tf/src/server/model"
	sess "shop.loadout.tf/src/server/session"
)

const minPasswordLen = 8
const maxPasswordLen = 72 // max bcrypt len

func apiCreateAccount(c *gin.Context, params map[string]any) apiError {
	if params == nil {
		return CreateApiError(NoParamsError)
	}

	username, ok := params["username"].(string)
	if !ok {
		return CreateApiError(InvalidParamUsername)
	}

	username = strings.ToLower(username)

	password, ok := params["password"].(string)
	if !ok {
		return CreateApiError(InvalidParamPassword)
	}

	if len(password) < minPasswordLen {
		return CreateApiError(InvalidParamPassword)
	}

	if len(password) > maxPasswordLen {
		return CreateApiError(InvalidParamPassword)
	}

	exist, err := shop.UsernameExist(username)
	if err != nil || exist {
		if err != nil {
			logger.Log(c, err)
		}
		return CreateApiError(UnexpectedError)
	}

	user, err := shop.CreateUser(username, password)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}
	log.Println(user)

	jsonSuccess(c, map[string]any{})

	return nil
}

func GetUser(username string, password string) (*model.User, error) {
	user, err := shop.FindUserByName(username, password)
	if err == shop.WrongPasswordError {
		return nil, fmt.Errorf("can't check user password %s %w", username, err)
	}
	if err != nil {
		return nil, fmt.Errorf("can't find user %s %w", username, err)
	}

	return user, nil
}

func apiLogin(c *gin.Context, s sessions.Session, params map[string]any) apiError {
	authSession := sess.GetAuthSession(c)
	if _, ok := authSession.Get("user_id").(string); ok {
		//return CreateApiError(AlreadyAuthenticated)
		authSession.Delete("user_id")
	}

	username, ok := params["username"].(string)
	if !ok {
		logger.Log(c, errors.New("username doen't exist in params"))
		return CreateApiError(InvalidParamUsername)
	}

	username = strings.ToLower(username)

	password, ok := params["password"].(string)
	if !ok {
		logger.Log(c, errors.New("password doen't exist in params"))
		return CreateApiError(InvalidParamPassword)
	}

	user, err := GetUser(username, password)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(AuthenticationError)
	}
	copySessionToUser(c, s, user.ID)

	authSession.Set("user_id", user.ID)
	if err := authSession.Save(); err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, map[string]any{
		"authenticated": true,
	})

	return nil
}

func apiLogout(c *gin.Context, s sessions.Session) apiError {
	copyUserToSession(c, s)

	if err := sess.RemoveAuthSession(c); err != nil {
		// Log the error but succeed
		logger.Log(c, err)
	}

	jsonSuccess(c, nil)

	return nil
}

func apiGetUser(c *gin.Context) apiError {
	authSession := sess.GetAuthSession(c)
	if userID, ok := authSession.Get("user_id").(string); ok {
		user, err := shop.FindUserByID(userID)
		if err != nil {
			logger.Log(c, err)
			return CreateApiError(NotAuthenticated)
		}
		jsonSuccess(c, map[string]any{
			"authenticated":  true,
			"display_name":   user.DisplayName,
			"email":          user.Email,
			"email_verified": user.EmailVerified,
			"currency":       user.Currency,
			"address":        user.Address,
		})
		return nil
	}

	return CreateApiError(NotAuthenticated)
}

func apiGetOrders(c *gin.Context) apiError {
	authSession := sess.GetAuthSession(c)
	if userID, ok := authSession.Get("user_id").(string); ok {
		user, err := shop.FindUserByID(userID)
		if err != nil {
			logger.Log(c, err)
			return CreateApiError(NotAuthenticated)
		}

		orders := make([]*model.Order, 0, len(user.Orders))
		for orderId := range user.Orders {
			order, err := shop.GetOrder(orderId)
			if err != nil {
				logger.Log(c, err)
			} else {
				orders = append(orders, order)
			}
		}

		jsonSuccess(c, map[string]any{
			"orders": orders,
		})
		return nil
	}

	return CreateApiError(UnexpectedError)
}

func copySessionToUser(c *gin.Context, s sessions.Session, userID string) error {
	// Copy favorites
	favorites, updateCurrency := s.Get("favorites").(map[string]any)
	if !updateCurrency {
		logger.Log(c, errors.New("favorites not found in session"))
	} else {
		shop.AddUserFavorites(userID, favorites)
	}

	// Copy cart
	cart, updateCurrency := s.Get("cart").(model.Cart)
	var updateCart bool
	if !updateCurrency {
		logger.Log(c, errors.New("cart not found in session"))
	} else {
		updateCart = cart.TotalQuantity() > 0
		/*
			if cart.TotalQuantity() > 0 {
				shop.SetUserCart(userID, cart)
			}
		*/
	}

	// Copy currency
	currency, updateCurrency := s.Get("currency").(string)
	/*
		if !currencyOk {
			currency = constants.DEFAULT_CURRENCY
		}
	*/
	//shop.SetUserCurrency(userID, currency)

	if updateCurrency || updateCart {
		err := shop.UpdateUser(model.User{ID: userID, Currency: currency, Cart: cart}, shop.UpdateUserFields{Currency: updateCurrency, Cart: updateCart})
		if err != nil {
			return err
		}
	}

	return nil
}

func copyUserToSession(c *gin.Context, s sessions.Session) error {
	authSession := sess.GetAuthSession(c)
	userID, ok := authSession.Get("user_id").(string)

	if !ok {
		return errors.New("invalid user_id")
	}

	user, err := shop.FindUserByID(userID)
	if err != nil {
		return fmt.Errorf("unable to find user %s: %w", userID, err)
	}

	// Copy currency
	s.Set("currency", user.Currency)

	// Copy favorites
	favorites := make(map[string]any)
	for favorite := range user.Favorites {
		favorites[favorite] = nil
	}
	s.Set("favorites", favorites)

	return nil
}

func apiSetUserInfos(c *gin.Context, params map[string]any) apiError {
	if params == nil {
		return CreateApiError(NoParamsError)
	}

	authSession := sess.GetAuthSession(c)
	userID, ok := authSession.Get("user_id").(string)

	if !ok {
		return CreateApiError(NotAuthenticated)
	}

	updateUserFields := shop.UpdateUserFields{}
	updateAny := false

	user, err := shop.FindUserByID(userID)
	if err != nil {
		logger.Log(c, fmt.Errorf("failed to get user in apiSetUserInfos %w", err))
		return CreateApiError(UnexpectedError)
	}

	if displayName, ok := params["display_name"].(string); ok && displayName != "" && user.DisplayName != displayName {
		updateUserFields.DisplayName = true
		user.DisplayName = displayName
		updateAny = true
	}

	if firstName, ok := params["address_first_name"].(string); ok && firstName != "" && user.Address.FirstName != firstName {
		updateUserFields.Address = true
		user.Address.FirstName = firstName
		updateAny = true
	}

	if lastName, ok := params["address_last_name"].(string); ok && lastName != "" && user.Address.LastName != lastName {
		updateUserFields.Address = true
		user.Address.LastName = lastName
		updateAny = true
	}

	if organization, ok := params["address_organization"].(string); ok && user.Address.Organization != organization {
		updateUserFields.Address = true
		user.Address.Organization = organization
		updateAny = true
	}

	if address1, ok := params["address_address1"].(string); ok && address1 != "" && user.Address.Address1 != address1 {
		updateUserFields.Address = true
		user.Address.Address1 = address1
		updateAny = true
	}

	if address2, ok := params["address_address2"].(string); ok && user.Address.Address2 != address2 {
		updateUserFields.Address = true
		user.Address.Address2 = address2
		updateAny = true
	}

	if city, ok := params["address_city"].(string); ok && city != "" && user.Address.City != city {
		updateUserFields.Address = true
		user.Address.City = city
		updateAny = true
	}

	if stateCode, ok := params["address_state_code"].(string); ok && user.Address.StateCode != stateCode {
		updateUserFields.Address = true
		user.Address.StateCode = stateCode
		updateAny = true
	}

	if stateName, ok := params["address_state_name"].(string); ok && user.Address.StateName != stateName {
		updateUserFields.Address = true
		user.Address.StateName = stateName
		updateAny = true
	}

	if countryCode, ok := params["address_country_code"].(string); ok && countryCode != "" && user.Address.CountryCode != countryCode {
		updateUserFields.Address = true
		user.Address.CountryCode = countryCode
		updateAny = true
	}

	if countryName, ok := params["address_country_name"].(string); ok && countryName != "" && user.Address.CountryName != countryName {
		updateUserFields.Address = true
		user.Address.CountryName = countryName
		updateAny = true
	}

	if postalCode, ok := params["address_postal_code"].(string); ok && postalCode != "" && user.Address.PostalCode != postalCode {
		updateUserFields.Address = true
		user.Address.PostalCode = postalCode
		updateAny = true
	}

	if phone, ok := params["address_phone"].(string); ok && user.Address.Phone != phone {
		updateUserFields.Address = true
		user.Address.Phone = phone
		updateAny = true
	}

	if email, ok := params["address_email"].(string); ok && email != "" && user.Address.Email != email {
		updateUserFields.Address = true
		user.Address.Email = email
		updateAny = true
	}

	if taxNumber, ok := params["address_tax_number"].(string); ok && user.Address.TaxNumber != taxNumber {
		updateUserFields.Address = true
		user.Address.TaxNumber = taxNumber
		updateAny = true
	}

	if updateAny {
		err := shop.UpdateUser(*user, updateUserFields)
		if err != nil {
			logger.Log(c, err)
			return CreateApiError(UnexpectedError)
		}
	}

	jsonSuccess(c, nil)
	return nil
}
