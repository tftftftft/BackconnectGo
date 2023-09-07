package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var proxyMap map[SrvInfo]BotInfo
var mapMutex sync.Mutex
var SrvIp, _ = GetPublicIP()

func AddProxy(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var newBotInfo BotInfo
	var newSrvInfo SrvInfo
	if err := json.Unmarshal(body, &newBotInfo); err != nil {
		log.Printf("Invalid JSON format: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	log.Printf("Received bot info: %+v", newBotInfo)

	// Check if the proxy already exists in the database
	exists, err := ProxyExistsInDatabase(newBotInfo.InitialProxyData.ProxyIP)
	if err != nil {
		log.Printf("Failed to check if proxy exists in database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if exists {
		log.Printf("Proxy already exists in database, not adding.")
		http.Error(w, "Proxy already exists", http.StatusBadRequest)
		return
	}

	randomPort := GenerateRandomPort()
	log.Printf("Generated random port for relay: %s", randomPort)

	newSrvInfo.ServerListeningPort = randomPort
	newSrvInfo.ServerIP = SrvIp

	mapMutex.Lock()

	proxyMap[newSrvInfo] = newBotInfo
	mapMutex.Unlock()
	log.Printf("Added bot with proxy: %+v", proxyMap[newSrvInfo])

	proxyAddress := net.JoinHostPort(newBotInfo.InitialProxyData.ProxyIP, newBotInfo.InitialProxyData.ProxyPort)
	go StartServer(randomPort, proxyAddress)

	log.Printf("UserID: %s", newBotInfo.UserID)

	if err := AddProxyToDatabase(newSrvInfo, newBotInfo); err != nil {
		log.Printf("Failed to add proxy to database: %v", err)
		http.Error(w, "Failed to add proxy to database", http.StatusInternalServerError)
		return
	}
	log.Printf("Proxy added to database")

	w.WriteHeader(http.StatusOK)
	log.Println("Successfully added proxy")
}

// function to return proxy info
func DisplayProxies(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to display proxies")
	mapMutex.Lock()
	defer mapMutex.Unlock()

	stringKeyedMap := make(map[string]BotInfo)

	for k, v := range proxyMap {
		keyStr, err := json.Marshal(k)
		if err != nil {
			log.Printf("Failed to serialize key: %v", err)
			continue
		}
		stringKeyedMap[string(keyStr)] = v
	}

	proxyJson, err := json.Marshal(stringKeyedMap)
	if err != nil {
		log.Printf("Failed to serialize data: %v", err)
		http.Error(w, "Failed to serialize data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(proxyJson)
	log.Printf("Proxies: %s", string(proxyJson))
}
func CheckProxiesAlive() {
	for {
		time.Sleep(5 * time.Minute) // Run every 5 minutes

		mapMutex.Lock()
		for srvInfo, botInfo := range proxyMap {
			proxyAddr := net.JoinHostPort(botInfo.InitialProxyData.ProxyIP, botInfo.InitialProxyData.ProxyPort)
			if !CheckProxyAlive(proxyAddr) {
				log.Printf("Proxy %s on port %s seems to be down", proxyAddr, srvInfo.ServerListeningPort)

				// Delete from database
				if err := DeleteProxyFromDatabase(botInfo.InitialProxyData.ProxyIP); err != nil {
					log.Printf("Failed to delete proxy from database: %v", err)
				}

				// Delete from map
				delete(proxyMap, srvInfo)
			}
		}
		mapMutex.Unlock()
	}
}
