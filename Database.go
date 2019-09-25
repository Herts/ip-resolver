package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func mysqlDb(conn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Println(err)
	}
	return db, err
}
