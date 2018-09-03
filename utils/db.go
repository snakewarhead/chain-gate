package utils

import (
	"database/sql"

	// self contain
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbPath = "./resources/chain_data.db"
	// dbPath = "./resources/chain_data_local.db"
)

var (
	DB *sql.DB
)

func init() {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	if err := DB.Ping(); err != nil {
		panic(err)
	}
}
