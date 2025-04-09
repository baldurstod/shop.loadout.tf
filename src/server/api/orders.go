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
		return errors.New("error while updating order")
	}

	err = createPrintfulOrder(order)
	if err != nil {
		log.Println(err)
		return errors.New("error while creating third party order")
	}

	return nil
}
