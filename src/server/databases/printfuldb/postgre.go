package printfuldb

import (
	"database/sql"

	_ "github.com/lib/pq"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/databases/postgre"
)

var printfulDb *sql.DB
var cacheMaxAge int64 = 86400

func InitPrintfulDB(config config.Database) {
	printfulDb = postgre.OpenPostgre(config.Datasource)
}
