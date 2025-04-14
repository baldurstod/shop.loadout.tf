package server

import (
	"log"
	"net/http"

	"shop.loadout.tf/src/server/databases"

	"github.com/gin-gonic/gin"
)

func imageHandler(c *gin.Context) {
	log.Println(c.FullPath(), c.Param("id"))

	img, err := databases.GetImage(c.Param("id"))
	if err != nil {
		c.String(http.StatusNotFound, "failed to read image")
		return
	}

	c.Data(http.StatusOK, "image/png", img)
}
