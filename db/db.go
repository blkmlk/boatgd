package db

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"log"
)

const DbFile = "./dbase.db"

var globalDb *db = nil

type db struct {
	dbase	*sql.DB
}

func GetDB() *db {
	var err error

	if globalDb == nil {
		globalDb = new(db)
		globalDb.dbase, err = sql.Open("sqlite3", DbFile)

		if err != nil {
			log.Fatal("DB:", err.Error())
		}

		globalDb.createDataTable()
	}

	return globalDb
}

func (dbase *db) createDataTable() {
	_, err := dbase.dbase.Exec("CREATE TABLE IF NOT EXISTS \"data\" (\"time\" INT," +
	"	\"pgn\" INT," +
	"	\"data\" BLOB)")

	if err != nil {
		log.Println(err.Error())
	}
}
