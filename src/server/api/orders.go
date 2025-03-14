package api

import (
	"errors"
	"log"

	"shop.loadout.tf/src/server/model"
	"shop.loadout.tf/src/server/mongo"
)

func approveOrder(order *model.Order) error {
	order.Status = "approved"
	err := mongo.UpdateOrder(order)
	if err != nil {
		log.Println(err)
		return errors.New("error while updating order in approveOrder")
	}

	err = createPrintfulOrder(order)
	if err != nil {
		log.Println(err)
		return errors.New("error while creting printful order")
	}

	return nil
}
