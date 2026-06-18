package shop

import (
	"database/sql"

	_ "github.com/lib/pq"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/databases/postgre"
)

var shopDb *sql.DB

func InitShopDB(config config.Database) {
	shopDb = postgre.OpenPostgre(config.Datasource)
}
