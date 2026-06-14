package printfuldb

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"shop.loadout.tf/src/server/config"
)

var printfulDb *sql.DB
var imagesDb *sql.DB

func InitPrintfulDB(config config.Database) {
	printfulDb = openPostgre(config.Datasource)
}

func InitImagesDB(config config.Database) {
	imagesDb = openPostgre(config.Datasource)
}

func openPostgre(dataSourceName string) *sql.DB {
	var err error
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	// db.Open() only creates a connection pool, and doesn't actually establish
	// a connection. To ensure the connection works you need to do *something*
	// with a connection.
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func ClosePostgre() {
	if printfulDb != nil {
		printfulDb.Close()
	}
	if imagesDb != nil {
		imagesDb.Close()
	}
}
