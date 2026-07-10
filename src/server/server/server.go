package server

import (
	"database/sql"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	assets "shop.loadout.tf"
	"shop.loadout.tf/src/server/api"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/databases/postgre"
	sess "shop.loadout.tf/src/server/session"
)

var ReleaseMode = "true"

var sessionsDb *sql.DB

func InitsessionsDB(config config.Database) {
	sessionsDb = postgre.OpenPostgre(config.Datasource)
}

func StartServer(config config.Config) {
	engine := initEngine(config)
	var err error

	log.Printf("Listening on port %d\n", config.HTTPS.Port)
	err = engine.RunTLS(":"+strconv.Itoa(config.HTTPS.Port), config.HttpsCertFile, config.HttpsKeyFile)
	log.Fatal(err)
}

func initEngine(config config.Config) *gin.Engine {
	if ReleaseMode == "true" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)

	for _, o := range config.AllowOrigins {
		if strings.Contains(o, "*") {
			panic("cors config must not have *")
		}
	}

	r.Use(cors.New(cors.Config{
		AllowMethods:    []string{"POST", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Request-Id"},
		AllowAllOrigins: false,
		AllowOrigins:    config.AllowOrigins,
		MaxAge:          12 * time.Hour,
	}))

	r.Use(secure.New(secure.Config{
		SSLRedirect:           true,
		STSSeconds:            315360000,
		ContentSecurityPolicy: "default-src 'self' *.paypal.com; img-src 'self' *.printful.com *.loadout.tf *.paypalobjects.com data:; object-src 'none'; frame-ancestors 'none'; script-src-elem 'self' 'unsafe-inline' *.paypal.com; style-src 'self' 'unsafe-inline';",
		ContentTypeNosniff:    true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
	}))

	var useFS fs.FS
	var assetsFs = &assets.Assets

	if ReleaseMode == "true" {
		fsys := fs.FS(assetsFs)
		useFS, _ = fs.Sub(fsys, "build/client")
	} else {
		useFS = os.DirFS("build/client")
	}

	// Init sessions store
	store, err := postgres.NewStore(sessionsDb, []byte(config.Sessions.Secret))
	if err != nil {
		log.Fatal(err)
	}

	r.Use(sessions.SessionsMany([]string{sess.RegularSession, sess.AuthSession}, store))
	r.Use(rewriteURL(r))
	r.StaticFS("/static", http.FS(useFS))
	r.POST("/api", api.ApiHandler)
	r.GET("/image/:id", imageHandler)

	return r
}

func rewriteURL(r *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api" {
			c.Next()
			return
		}
		if strings.HasPrefix(c.Request.URL.Path, "/image") {
			c.Next()
			return
		}
		if strings.HasPrefix(c.Request.URL.Path, "/@") {
			c.Request.URL.Path = "/"
			c.Abort()
			r.HandleContext(c)
			c.Next()
			return
		}
		if !strings.HasPrefix(c.Request.URL.Path, "/static") {
			c.Request.URL.Path = "/static" + c.Request.URL.Path
			c.Abort()
			r.HandleContext(c)
			c.Next()
			return
		}

		c.Next()
	}
}

/*
func initHandlers(config config.HTTP) *mux.Router {
	var assetsFs = &assets.Assets

	var useFS fs.FS

	if UseEmbed == "true" {
		fsys := fs.FS(assetsFs)
		useFS, _ = fs.Sub(fsys, "build/client")
	} else {
		useFS = os.DirFS("build/client")
	}

	r := mux.NewRouter()
	r.Use(rewriteURL)
	r.PathPrefix("/api").Handler(&RecoveryHandler{Handler: api.ApiHandler{}})
	r.PathPrefix("/").Handler(&RecoveryHandler{Handler: http.FileServer(http.FS(useFS))})

	return r
}

func rewriteURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/@") {
			r.URL.Path = "/"
		}
		next.ServeHTTP(w, r)
	})
}
*/
