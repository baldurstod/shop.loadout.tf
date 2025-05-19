package api

import (
	"errors"
	"fmt"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"shop.loadout.tf/src/server/databases"
	"shop.loadout.tf/src/server/logger"
	"shop.loadout.tf/src/server/model"
	sess "shop.loadout.tf/src/server/session"
)

const bcryptCost = 14
const minPasswordLen = 8
const maxPasswordLen = 72 // max bcrypt len

func apiCreateAccount(c *gin.Context, s sessions.Session, params map[string]any) apiError {
	if params == nil {
		return CreateApiError(NoParamsError)
	}

	username, ok := params["username"].(string)
	if !ok {
		return CreateApiError(InvalidParamUsername)
	}

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

	exist, err := databases.UsernameExist(username)
	if err != nil || exist {
		return CreateApiError(UnexpectedError)
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	user, err := databases.CreateUser(username, hashedPassword)
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}
	log.Println(user)

	jsonSuccess(c, map[string]any{})

	return nil
}

func verifyEmail(user *model.User) error {
	if user.EmailVerified {
		return nil
	}

	return nil

}

func GetUser(username string, password string) (*model.User, error) {
	user, err := databases.FindUserByName(username)
	if err != nil {
		return nil, fmt.Errorf("can't find user %s", username)
	}

	if user.Password == "" {
		return nil, fmt.Errorf("user %s has an empty password", username)
	}

	if !CheckPasswordHash(password, user.Password) {
		return nil, errors.New("wrong password")
	}

	return user, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func apiLogin(c *gin.Context, s sessions.Session, params map[string]any) apiError {
	username, ok := params["username"].(string)
	if !ok {
		return CreateApiError(InvalidParamUsername)
	}

	password, ok := params["password"].(string)
	if !ok {
		return CreateApiError(InvalidParamPassword)
	}

	user, err := GetUser(username, password)
	if err != nil {
		return CreateApiError(AuthenticationError)
	}
	copySessionToUser(c, s, user.ID)

	authSession := sess.GetAuthSession(c)
	authSession.Set("user_id", user.ID)
	if err := authSession.Save(); err != nil {
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, map[string]any{})

	return nil
}

func apiLogout(c *gin.Context, s sessions.Session) apiError {
	copyUserToSession(c, s)

	if err := sess.RemoveAuthSession(c); err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, map[string]any{})

	return nil
}

func apiGetuser(c *gin.Context, s sessions.Session) apiError {
	authSession := sess.GetAuthSession(c)
	if userID, ok := authSession.Get("user_id").(string); ok {
		user, err := databases.FindUserByID(userID)
		if err != nil {
			logger.Log(c, err)
			jsonSuccess(c, map[string]any{"authenticated": false})
		}
		jsonSuccess(c, map[string]any{
			"authenticated": true,
			"display_name":  user.DisplayName,
		})
		return nil
	}

	jsonSuccess(c, map[string]any{"authenticated": false})
	return nil
}

func copySessionToUser(c *gin.Context, s sessions.Session, userID string) error {
	// Copy favorites
	favorites, ok := s.Get("favorites").(map[string]any)
	if !ok {
		logger.Log(c, errors.New("favorites not found in session"))
	} else {
		databases.AddUserFavorites(userID, favorites)
	}

	// Copy cart
	cart, ok := s.Get("cart").(model.Cart)
	if !ok {
		logger.Log(c, errors.New("cart not found in session"))
	} else {
		if cart.TotalQuantity() > 0 {
			databases.SetUserCart(userID, cart)
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

	user, err := databases.FindUserByID(userID)
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
