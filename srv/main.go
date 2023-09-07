package main

import (
	"log"
	"net/http"
)

func main() {

	//init random seed
	InitSeed()
	ConnectDB()
	//create proxy map
	proxyMap = make(map[SrvInfo]BotInfo)

	//routes
	http.HandleFunc("/addProxy", AddProxy)
	http.HandleFunc("/proxies", DisplayProxies)
	log.Println("Starting proxy server")

	//create http server to handle routes
	go func() {
		if err := http.ListenAndServe("0.0.0.0:20001", nil); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()
	log.Println("Listening on port 20001")

	//for each new proxy starg relaying traffic
	for SrvInfo, BotInfo := range proxyMap {
		go StartServer(SrvInfo.ServerListeningPort, BotInfo.ProxyIP)
	}

	//goroutine to check proxies are alive periodically
	go CheckProxiesAlive()

	//keep main goroutine alive
	select {}
}
