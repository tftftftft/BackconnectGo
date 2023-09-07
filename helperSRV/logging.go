package main

import (
	"log"
	"os"
)

func SetupLogging() {
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}
	log.SetOutput(logFile)
}
