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

func apiSendCurrentEmailVerification(c *gin.Context) apiError {
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

	if err = verifyEmail(userID, user.Email, verifyCurrentEmail); err != nil {
		logger.Log(c, fmt.Errorf("failed to send mail verification <%w>", err))
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, nil)
	return nil
}

func apiSendNewEmailVerification(c *gin.Context, params map[string]any) apiError {
	if params == nil {
		return CreateApiError(NoParamsError)
	}

	authSession := sess.GetAuthSession(c)
	userId, ok := authSession.Get("user_id").(string)

	if !ok {
		return CreateApiError(NotAuthenticated)
	}

	var err error

	var email string
	if email, ok = params["email"].(string); !ok {
		return CreateApiError(InvalidParamEmail)
	}

	// Check email validity
	if _, err := mail.ParseAddress(email); err != nil {
		return CreateApiError(InvalidParamEmail)
	}

	var user *model.User
	if user, err = shop.FindUserByID(userId); err != nil {
		logger.Log(c, fmt.Errorf("failed to get user in apiSendNewEmailVerification <%w>", err))
		return CreateApiError(UnexpectedError)
	}

	// Don't do anything if mails are identical
	if user.Email == email {
		return CreateApiError(InvalidParamEmail)
	}

	if err = verifyEmail(userId, email, verifyNewEmail); err != nil {
		logger.Log(c, fmt.Errorf("failed to send mail verification <%w>", err))
		return CreateApiError(UnexpectedError)
	}

	jsonSuccess(c, nil)
	return nil
}

func apiVerifyCurrentEmail(c *gin.Context, params map[string]any) apiError {
	if params == nil {
		return CreateApiError(NoParamsError)
	}

	authSession := sess.GetAuthSession(c)
	userId, ok := authSession.Get("user_id").(string)

	if !ok {
		return CreateApiError(NotAuthenticated)
	}

	var user *model.User
	var err error

	if user, err = shop.FindUserByID(userId); err != nil {
		logger.Log(c, fmt.Errorf("failed to get user in apiVerifyEmail <%w>", err))
		return CreateApiError(UnexpectedError)
	}

	var code string
	if code, ok = params["code"].(string); !ok {
		return CreateApiError(InvalidParamCode)
	}

	codeUserId, err := shop.CheckEmailVerification(code, user.Email)
	if err != nil {
		if err == shop.ErrCodeExpired {
			return CreateApiError(ExpiredVerificationCode)
		} else {
			return CreateApiError(InvalidVerificationCode)
		}
	}

	if codeUserId != userId {
		return CreateApiError(InvalidVerificationCode)
	}

	jsonSuccess(c, nil)
	return nil
}

func apiVerifyNewEmail(c *gin.Context, params map[string]any) apiError {
	if params == nil {
		return CreateApiError(NoParamsError)
	}

	authSession := sess.GetAuthSession(c)
	userId, ok := authSession.Get("user_id").(string)

	if !ok {
		return CreateApiError(NotAuthenticated)
	}

	var email string
	if email, ok = params["email"].(string); !ok {
		return CreateApiError(InvalidParamEmail)
	}

	// Check email validity
	if _, err := mail.ParseAddress(email); err != nil {
		return CreateApiError(InvalidParamEmail)
	}

	var code string
	if code, ok = params["code"].(string); !ok {
		return CreateApiError(InvalidParamCode)
	}

	codeUserId, err := shop.CheckEmailVerification(code, email)
	if err != nil {
		if err == shop.ErrCodeExpired {
			return CreateApiError(ExpiredVerificationCode)
		} else {
			return CreateApiError(InvalidVerificationCode)
		}
	}

	if codeUserId != userId {
		return CreateApiError(InvalidVerificationCode)
	}

	jsonSuccess(c, nil)
	return nil
}

func apiChangeEmail(c *gin.Context, params map[string]any) apiError {
	if params == nil {
		return CreateApiError(NoParamsError)
	}

	authSession := sess.GetAuthSession(c)
	userId, ok := authSession.Get("user_id").(string)

	if !ok {
		return CreateApiError(NotAuthenticated)
	}

	var currentCode string
	if currentCode, ok = params["current_code"].(string); !ok {
		return CreateApiError(InvalidParamCode)
	}

	var newEmail string
	if newEmail, ok = params["new_email"].(string); !ok {
		return CreateApiError(InvalidParamEmail)
	}

	// Check email validity
	if _, err := mail.ParseAddress(newEmail); err != nil {
		return CreateApiError(InvalidParamEmail)
	}

	var newCode string
	if newCode, ok = params["new_code"].(string); !ok {
		return CreateApiError(InvalidParamCode)
	}

	var user *model.User
	var err error

	if user, err = shop.FindUserByID(userId); err != nil {
		logger.Log(c, fmt.Errorf("failed to get user in apiChangeEmail <%w>", err))
		return CreateApiError(UnexpectedError)
	}

	// Don't do anything if mails are identical
	if user.Email == newEmail {
		return CreateApiError(UnexpectedError)
	}

	var currentCodeUserId string
	if user.Email != "" {
		currentCodeUserId, err = shop.CheckEmailVerification(currentCode, user.Email)
		if err != nil {
			if err == shop.ErrCodeExpired {
				return CreateApiError(ExpiredVerificationCode)
			} else {
				return CreateApiError(InvalidVerificationCode)
			}
		}
	}

	newCodeUserId, err := shop.CheckEmailVerification(newCode, newEmail)
	if err != nil {
		if err == shop.ErrCodeExpired {
			return CreateApiError(ExpiredVerificationCode)
		} else {
			return CreateApiError(InvalidVerificationCode)
		}
	}

	if user.Email == "" {
		// Only check the new code
		if newCodeUserId != userId {
			return CreateApiError(InvalidVerificationCode)
		}
	} else {
		// Check old and new code
		if currentCodeUserId != userId || newCodeUserId != userId {
			return CreateApiError(InvalidVerificationCode)
		}
	}

	// We checked old code and new code, update the email
	user.Email = newEmail
	user.EmailVerified = true
	err = shop.UpdateUser(*user, shop.UpdateUserFields{Email: true, EmailVerified: true})
	if err != nil {
		logger.Log(c, err)
		return CreateApiError(UnexpectedError)
	}

	// Delete codes
	if err = shop.DeleteEmailVerification(userId); err != nil {
		logger.Log(c, fmt.Errorf("failed to purge verification codes for user %s in verifyEmailHandler <%w>", userId, err))
	}

	jsonSuccess(c, nil)
	return nil
}

func verifyEmail(userId string, mail string, text string) error {
	code := randstr.Hex(32)

	if err := shop.InsertEmailVerification(userId, mail, code); err != nil {
		return err
	}

	if err := email.SendMailVerification(mail, code, text); err != nil {
		return err
	}

	return nil
}

var verifyCurrentEmail = `
	<html>
	<body>
	<h1>Verify your email address</h1>
	To change your email address, we need to make sure you have access to this email address.<br>
	To verify your email address, use the following code:
	<h1>{{.code}}</h1>
	Don't share this code with anyone. Never use this code outside the official website.<br>
	<br>
	If you didn't request this code nor trying to change your email address, your account may be compromised and you must change your password.<br>
	<br>
	This code replace any previous code sent to this email address.<br>
	</body>
	</html>

	`

var verifyNewEmail = `
	<html>
	<body>
	<h1>Verify your email address</h1>
	To change your email address, we need to make sure you have access to this email address.<br>
	To verify your email address, use the following code:
	<h1>{{.code}}</h1>
	Don't share this code with anyone. Never use this code outside the official website.<br>
	<br>
	If you didn't request this code, you can safely ignore this email.<br>
	<br>
	This code replace any previous code sent to this email address.<br>
	</body>
	</html>
	`
