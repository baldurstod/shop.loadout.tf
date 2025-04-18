package api

import (
	"fmt"

	"shop.loadout.tf/src/server/databases"
	"shop.loadout.tf/src/server/model"
)

func approveOrder(order *model.Order) error {
	order.Status = "approved"
	err := databases.UpdateOrder(order)
	if err != nil {
		return fmt.Errorf("error while updating order: %w", err)
	}

	err = createPrintfulOrder(order)
	if err != nil {
		return fmt.Errorf("error while creating printful order: %w", err)
	}

	return nil
}
