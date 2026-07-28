package server

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"shop.loadout.tf/src/server/databases/shop"
	"shop.loadout.tf/src/server/logger"
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
		if err == shop.ErrCodeExpired {
			c.String(http.StatusOK, "This code is no longer valid.")
		} else {
			c.String(http.StatusInternalServerError, "error: wrong or expired code")
		}
		return
	}

	user, err := shop.FindUserByID(userId)
	if err != nil {
		logger.Log(c, fmt.Errorf("failed to get user %s in verifyEmailHandler <%w>", userId, err))
		c.String(http.StatusInternalServerError, "error: wrong or expired code")
		return
	}

	if user.Email != email {
		logger.Log(c, fmt.Errorf("user email %s doesn't match verification email %s for user %s", user.Email, email, userId))
		c.String(http.StatusInternalServerError, "error: wrong or expired code")
		return
	}

	if err = shop.UpdateUser(model.User{ID: userId, EmailVerified: true}, shop.UpdateUserFields{EmailVerified: true}); err != nil {
		logger.Log(c, fmt.Errorf("failed to update user %s in verifyEmailHandler <%w>", userId, err))
		c.String(http.StatusInternalServerError, "unexpected error: contact support")
		return
	}

	if err = shop.DeleteEmailVerification(userId); err != nil {
		logger.Log(c, fmt.Errorf("failed to purge verification codes for user %s in verifyEmailHandler <%w>", userId, err))
	}

	c.Redirect(http.StatusFound, "/@user")
}
