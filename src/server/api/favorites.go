package api

import (
	"errors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"shop.loadout.tf/src/server/databases"
	"shop.loadout.tf/src/server/logger"
	"shop.loadout.tf/src/server/model"
)

func apiGetFavorites(c *gin.Context, s sessions.Session) apiError {
	if userID, ok := s.Get("user_id").(string); ok {
		user, err := databases.FindUserByID(userID)
		if err != nil {
			logger.Log(c, err)
			return CreateApiError(UnexpectedError)
		}
		jsonSuccess(c, map[string]any{"favorites": apiGetUserFavorites(user)})
		return nil
	}

	jsonSuccess(c, map[string]any{"favorites": apiGetSessionFavorites(s)})
	return nil
}

func apiGetSessionFavorites(s sessions.Session) []string {
	favorites, ok := s.Get("favorites").(map[string]any)
	if !ok {
		favorites = make(map[string]any)
	}

	v := make([]string, 0, len(favorites))

	for key := range favorites {
		v = append(v, key)
	}
	return v
}

func apiGetUserFavorites(user *model.User) []string {
	v := make([]string, 0, len(user.Favorites))

	for key := range user.Favorites {
		v = append(v, key)
	}
	return v
}

func apiSetFavorite(c *gin.Context, s sessions.Session, params map[string]any) apiError {
	if params == nil {
		return CreateApiError(NoParamsError)
	}

	productID, ok := params["product_id"].(string)
	if !ok {
		return CreateApiError(InvalidParamProductID)
	}

	isFavorite, ok := params["is_favorite"].(bool)
	if !ok {
		return CreateApiError(InvalidParamIsFavorite)
	}

	if userID, ok := s.Get("user_id").(string); ok {
		user, err := databases.FindUserByID(userID)
		if err != nil {
			logger.Log(c, err)
			return CreateApiError(UnexpectedError)
		}

		if apiSetUserFavorite(user, productID, isFavorite) != nil {
			return CreateApiError(UnexpectedError)
		}

		return nil
	}

	jsonSuccess(c, nil)
	return nil
}

func apiSetSessionFavorite(s sessions.Session, productID string, isFavorite bool) error {
	favorites, ok := s.Get("favorites").(map[string]any)
	if !ok {
		return errors.New("favorites not found")
	}

	if isFavorite {
		favorites[productID] = struct{}{}
	} else {
		delete(favorites, productID)
	}
	return nil
}

func apiSetUserFavorite(user *model.User, productID string, isFavorite bool) error {
	/*
		favorites, ok := s.Get("favorites").(map[string]any)
		if !ok {
			return errors.New("favorites not found")
		}

		if isFavorite {
			favorites[productId] = struct{}{}
		} else {
			delete(favorites, productId)
		}
	*/
	return nil
}
