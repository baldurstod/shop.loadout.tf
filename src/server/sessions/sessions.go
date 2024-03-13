package sessions

import (
	"encoding/base64"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"shop.loadout.tf/src/server/config"
)

var store *sessions.FilesystemStore

func InitSessions(config config.Sessions) {
	authKey, err := base64.StdEncoding.DecodeString(config.AuthKey)
	if err != nil {
		log.Fatal(err)
	}

	encryptKey, err := base64.StdEncoding.DecodeString(config.EncryptKey)
	if err != nil {
		log.Fatal(err)
	}

	store = sessions.NewFilesystemStore(config.Path, authKey, encryptKey)
	store.MaxLength(50000)

	log.Println(store)
}

func GetSession(r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "session_id")
	return session
}
