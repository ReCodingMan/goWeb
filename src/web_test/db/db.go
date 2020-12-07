package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func init()  {
	var err error
	Db, err = sql.Open("mysql", "root:root@/test_db?charset=utf8mb4&parseTime=true")
	if err != nil {
		panic(err)
	}
}