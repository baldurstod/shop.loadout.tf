package server

import (
	"github.com/gorilla/mux"
	"io/fs"
	"log"
	"net/http"
	"os"
	"shop.loadout.tf"
	"shop.loadout.tf/src/server/api"
	"shop.loadout.tf/src/server/config"
	"strconv"
	"strings"
)

var UseEmbed = "true"

func StartServer(config config.HTTP) {
	handler := initHandlers(config)

	log.Printf("Listening on port %d\n", config.Port)
	err := http.ListenAndServeTLS(":"+strconv.Itoa(config.Port), config.HttpsCertFile, config.HttpsKeyFile, handler)
	log.Fatal(err)
}

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
