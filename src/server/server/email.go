package server

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"shop.loadout.tf/src/server/databases/shop"
	"shop.loadout.tf/src/server/model"
)

func verifyEmailHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.String(http.StatusInternalServerError, "error")
			log.Println(err, string(debug.Stack()))
		}
	}()

	code := c.Query("code")
	email := c.Query("email")

	userId, err := shop.CheckEmailVerification(code, email)
	if err != nil {
		if err == shop.ErrCodeValidity {
			c.String(http.StatusOK, "This code is no longer valid.")
		} else {
			c.String(http.StatusInternalServerError, "error")
		}
		return
	}
	err = shop.UpdateUser(model.User{ID: userId, EmailVerified: true}, shop.UpdateUserFields{EmailVerified: true})
	if err != nil {
		return
	}

	c.Redirect(http.StatusFound, "/@user")
}
