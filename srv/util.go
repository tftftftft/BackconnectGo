package main

import (
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func InitSeed() {
	rand.Seed(time.Now().UnixNano())
}

// function to handle health check requests
func CheckProxyAlive(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	// SOCKS5 handshake is [0x05, 0x01, 0x00]
	if _, err = conn.Write([]byte{0x05, 0x01, 0x00}); err != nil {
		return false
	}

	response := make([]byte, 2)
	if _, err = io.ReadFull(conn, response); err != nil {
		return false
	}

	// Expected response is [0x05, 0x00]
	if response[0] != 0x05 || response[1] != 0x00 {
		return false
	}

	return true
}

// function to
func GenerateRandomPort() string {
	for {
		candidate := rand.Intn(10001) + 40000 // Generate a number between 40000 and 50000
		candidateStr := strconv.Itoa(candidate)
		log.Printf("Generated candidate port: %s", candidateStr)

		mapMutex.Lock()
		candidateSrvInfo := SrvInfo{ServerIP: SrvIp, ServerListeningPort: candidateStr}
		_, exists := proxyMap[candidateSrvInfo]
		mapMutex.Unlock()
		if exists {
			log.Printf("Port %s already exists in proxyMap, generating a new one", candidateStr)
			continue
		}

		listener, err := net.Listen("tcp", "0.0.0.0:"+candidateStr)
		if err != nil {
			log.Printf("Port %s is already in use by another application, generating a new one", candidateStr)
			continue
		}
		listener.Close()
		log.Printf("Port %s is available", candidateStr)
		return candidateStr
	}
}

func GetPublicIP() (string, error) {
	res, err := http.Get("https://api4.ipext.org")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	ip := string(bodyBytes)
	return ip, nil
}

func GenerateUUID() string {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		log.Println("Error generating UUID:", err)
		return ""
	}
	return newUUID.String()
}
