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
