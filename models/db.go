package models

import (
	"database/sql"
	"reflect"

	// self contain
	_ "github.com/go-sql-driver/mysql"
)

const (
	// change it in proction
	dbPath = "root:123456@tcp(192.168.1.2:3306)/divide?charset=utf8"
	dbMaxOpenConns = 10
	dbMaxIdelConns = 10
)

var (
	db *sql.DB
)

func init() {
	var err error
	db, err = sql.Open("mysql", dbPath)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(dbMaxOpenConns)
	db.SetMaxIdleConns(dbMaxIdelConns)
}

// list the pointer of feilds of the struct of the model struct, usually use in the row.scan
// note that the model must be a pointer to a struct
func dbColumns(model interface{}) []interface{} {
	v := reflect.ValueOf(model).Elem()

	num := v.NumField()
	fields := make([]interface{}, num, num)
	for i := 0; i < num; i++ {
		fields[i] = v.Field(i).Addr().Interface()
	}

	return fields
}
