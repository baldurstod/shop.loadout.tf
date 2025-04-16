package api

import (
	"errors"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func apiGetFavorites(c *gin.Context, s sessions.Session) error {
	favorites, ok := s.Get("favorites").(map[string]any)
	if !ok {
		favorites = make(map[string]any)
	}

	v := make([]string, 0, len(favorites))

	for key := range favorites {
		v = append(v, key)
	}

	jsonSuccess(c, map[string]any{"favorites": v})
	return nil
}

func apiSetFavorite(c *gin.Context, s sessions.Session, params map[string]any) error {
	if params == nil {
		return errors.New("no params provided")
	}

	productId, ok := params["product_id"].(string)
	if !ok {
		return errors.New("missing params product_id")
	}

	isFavorite, ok := params["is_favorite"].(bool)
	if !ok {
		return errors.New("missing params is_favorite")
	}

	favorites, ok := s.Get("favorites").(map[string]any)
	if !ok {
		return errors.New("favorites not found")
	}

	if isFavorite {
		favorites[productId] = struct{}{}
	} else {
		delete(favorites, productId)
	}

	log.Println(favorites)

	jsonSuccess(c, nil)
	return nil
}
