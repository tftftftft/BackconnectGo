package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// SendErrorResponse sends an HTTP error response
func SendErrorResponse(w http.ResponseWriter, msg string, code int) {
	http.Error(w, msg, code)
}

func newProxyRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("Received new proxy registration request.")

	body, err := ReadRequestBody(r)
	if err != nil {
		log.Printf("Failed reading request body: %v", err)
		SendErrorResponse(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var initialData InitialProxyData
	if err := json.Unmarshal(body, &initialData); err != nil {
		log.Printf("Failed unmarshalling initial proxy data: %v", err)
		SendErrorResponse(w, "Failed to unmarshal initial proxy data", http.StatusBadRequest)
		return
	}

	IPInfo, err := FetchIPInfo(initialData.ProxyIP)
	if err != nil {
		log.Printf("Failed to fetch IP info: %v", err)
		SendErrorResponse(w, "Failed to fetch IP info", http.StatusInternalServerError)
		return
	}

	BotInfo := BotInfo{
		InitialProxyData: initialData,
		IPInfo:           *IPInfo,
	}

	if BotInfoBody, err := json.Marshal(BotInfo); err == nil {
		resp, err := ForwardRequest(BotInfoBody)
		if err != nil {
			log.Printf("Failed to forward request: %v", err)
			SendErrorResponse(w, "Failed to forward request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		log.Printf("Received response status code: %d", resp.StatusCode)
	} else {
		log.Printf("Failed to marshal bot info: %v", err)
		SendErrorResponse(w, "Failed to marshal bot info", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Successfully sent OK status to original client.")
}

func main() {
	SetupLogging()
	log.Println("Starting the server on port 25000.")

	http.HandleFunc("/route", newProxyRegister)
	if err := http.ListenAndServe(":25000", nil); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
