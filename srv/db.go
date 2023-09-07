package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func ConnectDB() {
	var err error
	db, err = sql.Open("mysql", "user:pass@tcp(ip:3306)/name")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
}

func AddProxyToDatabase(srvInfo SrvInfo, botInfo BotInfo) error {
	// SQL query to insert data into Proxies table
	query := `INSERT INTO Proxies (
        ServerIP,
        ServerListeningPort,
        ProxyIP,
        ProxyPort,
        UserID,
        BuildVersion,
        Continent,
        ContinentCode,
        Country,
        CountryCode,
        Region,
        RegionName,
        City,
        Zip,
        Timezone,
        ISP,
        Org,
        ASName,
        Mobile,
        Proxy,
        Hosting
    ) VALUES (
        ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
    )`

	// Prepare the SQL query
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL query
	_, err = stmt.Exec(
		srvInfo.ServerIP,
		srvInfo.ServerListeningPort,
		botInfo.ProxyIP,
		botInfo.ProxyPort,
		botInfo.UserID,
		botInfo.BuildVersion,
		botInfo.Continent,
		botInfo.ContinentCode,
		botInfo.Country,
		botInfo.CountryCode,
		botInfo.Region,
		botInfo.RegionName,
		botInfo.City,
		botInfo.Zip,
		botInfo.Timezone,
		botInfo.ISP,
		botInfo.Org,
		botInfo.ASName,
		botInfo.Mobile,
		botInfo.Proxy,
		botInfo.Hosting,
	)
	if err != nil {
		return err
	}

	return nil
}

func ProxyExistsInDatabase(proxyIP string) (bool, error) {
	// SQL query to check if proxy exists
	checkQuery := `SELECT COUNT(*) FROM Proxies WHERE ProxyIP = ?`

	var count int
	err := db.QueryRow(checkQuery, proxyIP).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func DeleteProxyFromDatabase(proxyIP string) error {
	deleteQuery := "DELETE FROM Proxies WHERE ProxyIP = ?"
	_, err := db.Exec(deleteQuery, proxyIP)
	return err
}
