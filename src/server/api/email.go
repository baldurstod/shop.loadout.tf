package api

import (
	"fmt"
	"net/mail"

	"github.com/baldurstod/randstr"
	"github.com/gin-gonic/gin"
	"shop.loadout.tf/src/server/databases/shop"
	"shop.loadout.tf/src/server/email"
	"shop.loadout.tf/src/server/logger"
	"shop.loadout.tf/src/server/model"
	sess "shop.loadout.tf/src/server/session"
)

func apiSendEmailVerification(c *gin.Context, params map[string]any) apiError {
	if params == nil {
		return CreateApiError(NoParamsError)
	}

	authSession := sess.GetAuthSession(c)
	userID, ok := authSession.Get("user_id").(string)

	if !ok {
		return CreateApiError(NotAuthenticated)
	}

	var user *model.User
	var err error

	if user, err = shop.FindUserByID(userID); err != nil {
		logger.Log(c, fmt.Errorf("failed to get user in apiVerifyEmail <%w>", err))
		return CreateApiError(UnexpectedError)
	}

	if user.EmailVerified {
		return CreateApiError(UnexpectedError)
	}

	var email string
	if email, ok = params["email"].(string); !ok {
		return CreateApiError(InvalidParamEmail)
	}

	// Check email validity
	if _, err := mail.ParseAddress(email); err != nil {
		return CreateApiError(InvalidParamEmail)
	}

	if err = verifyEmail(userID, email); err != nil {
		logger.Log(c, fmt.Errorf("failed to send mail verification <%w>", err))
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, nil)
	return nil
}

func verifyEmail(userID string, mail string) error {
	code := randstr.Hex(32)

	if err := shop.InsertEmailVerification(userID, mail, code); err != nil {
		return err
	}

	if err := email.SendMailVerification(mail, code); err != nil {
		return err
	}

	return nil
}
