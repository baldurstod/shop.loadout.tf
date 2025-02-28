package server

import (
	"log"
	"net/http"

	"shop.loadout.tf/src/server/mongo"

	"github.com/gin-gonic/gin"
)

func imageHandler(c *gin.Context) {
	log.Println(c.FullPath(), c.Param("id"))

	img, err := mongo.GetImage(c.Param("id"))
	if err != nil {
		c.String(http.StatusNotFound, "failed to read image")
		return
	}

	c.Data(http.StatusOK, "image/png", img)
}
