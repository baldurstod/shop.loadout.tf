package api

import (
	"errors"
	"log"

	"shop.loadout.tf/src/server/databases"
	"shop.loadout.tf/src/server/model"
)

func approveOrder(order *model.Order) error {
	order.Status = "approved"
	err := databases.UpdateOrder(order)
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
