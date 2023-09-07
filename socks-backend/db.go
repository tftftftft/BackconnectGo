package main

import (
	"database/sql"
	"log"
)

func ConnectDB() {
	var err error
	db, err = sql.Open("mysql", "proxyBase:DmAryiDjdmfS3Jt2@tcp(38.180.61.247:3306)/proxyBase")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
}
