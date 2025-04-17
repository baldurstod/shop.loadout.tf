package server

import (
	"context"
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
	"github.com/gin-contrib/sessions/mongo/mongodriver"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	assets "shop.loadout.tf"
	"shop.loadout.tf/src/server/api"
	"shop.loadout.tf/src/server/config"
	sess "shop.loadout.tf/src/server/session"
)

var ReleaseMode = "true"

func StartServer(config config.Config) {
	engine := initEngine(config)
	var err error

	log.Printf("Listening on port %d\n", config.Port)
	err = engine.RunTLS(":"+strconv.Itoa(config.Port), config.HttpsCertFile, config.HttpsKeyFile)
	log.Fatal(err)
}

func initEngine(config config.Config) *gin.Engine {
	if ReleaseMode == "true" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)

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
		FrameDeny:             true,
		ContentSecurityPolicy: "default-src 'self'; img-src 'self' *.printful.com *.loadout.tf; object-src 'none'",
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
	mongoOptions := options.Client().ApplyURI(config.Sessions.DB.ConnectURI)
	client, err := mongo.Connect(context.Background(), mongoOptions)
	if err != nil {
		log.Fatal(err)
	}
	c := client.Database(config.Sessions.DB.DBName).Collection("sessions")
	store := mongodriver.NewStore(c, 86400*30, true, []byte(config.Sessions.Secret))

	r.Use(sessions.Sessions(config.Sessions.SessionName, store))
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
			r.HandleContext(c)
			c.Next()
			return
		}
		if !strings.HasPrefix(c.Request.URL.Path, "/static") {
			c.Request.URL.Path = "/static" + c.Request.URL.Path
			r.HandleContext(c)
			c.Next()
			return
		}
		session := sess.GetSession(c)
		sess.SaveSession(session)

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
