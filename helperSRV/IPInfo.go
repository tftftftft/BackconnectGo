package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func FetchIPInfo(ip string) (*IPInfo, error) {
	resp, err := http.Get("http://ip-api.com/json/" + ip + "?fields=status,message,continent,continentCode,country,countryCode,region,regionName,city,zip,timezone,isp,org,as,asname,mobile,proxy,hosting")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info IPInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}
	log.Println(info)
	return &info, nil
}
