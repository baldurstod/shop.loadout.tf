package server

import (
	"log"
	"net/http"
	"runtime/debug"
)

type RecoveryHandler struct {
	Handler http.Handler
}

func (h *RecoveryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err, string(debug.Stack()))
		}
	}()

	h.Handler.ServeHTTP(w, r)
}
