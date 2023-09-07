package main

import (
	"database/sql"
	"log"
)

func ConnectDB() {
	var err error
	db, err = sql.Open("mysql", "user:pswd@tcp(ip:3306)/base")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
}
