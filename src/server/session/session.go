package sess

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var RegularSession = "session1"
var AuthSession = "session2"

func GetRegularSession(c *gin.Context) sessions.Session {
	session := sessions.DefaultMany(c, RegularSession)
	log.Println(session.ID())
	session.Options(sessions.Options{MaxAge: 86400 * 30, Path: "/", Secure: true, HttpOnly: true, SameSite: http.SameSiteStrictMode})
	return session
}

func GetAuthSession(c *gin.Context) sessions.Session {
	session := sessions.DefaultMany(c, AuthSession)
	log.Println(session.ID())
	session.Options(sessions.Options{MaxAge: 86400 * 5, Path: "/", Secure: true, HttpOnly: true, SameSite: http.SameSiteStrictMode})
	return session
}

func RemoveAuthSession(c *gin.Context) error {
	session := sessions.DefaultMany(c, AuthSession)
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1})
	return session.Save()
}
