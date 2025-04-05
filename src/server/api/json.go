package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
func writeJSON(w *http.ResponseWriter, r *http.Request, datas *map[string]any) {
	(*w).Header().Add("Content-Type", "application/json")
	(*w).Header().Add("Access-Control-Allow-Origin", "*")
	(*w).Header().Add("Access-Control-Allow-Headers", "content-type")

	if datas != nil {
		j, err := json.Marshal(datas)
		if err == nil {
			(*w).Write(j)
			return
		}
	}
}
*/

func jsonError(c *gin.Context, e error) {
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"error":   e.Error(),
	})
}

func jsonSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  data,
	})
}
