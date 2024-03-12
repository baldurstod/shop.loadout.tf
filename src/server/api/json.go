package api

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w *http.ResponseWriter, r *http.Request, datas *map[string]interface{}) {
	(*w).Header().Add("Content-Type", "application/json")
	(*w).Header().Add("Access-Control-Allow-Origin", "*")

	if datas != nil {
		j, err := json.Marshal(datas)
		if err == nil {
			(*w).Write(j)
			return
		}
	}
}

func jsonError(w http.ResponseWriter, r *http.Request, e error) {
	failure := map[string]interface{}{
		"success": false,
		"error":   e.Error(),
	}
	writeJSON(&w, r, &failure)
}

func jsonSuccess(w http.ResponseWriter, r *http.Request, data interface{}) {
	failure := map[string]interface{}{
		"success": true,
		"result":  data,
	}
	writeJSON(&w, r, &failure)
}
