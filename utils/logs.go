package utils

import (
	"log"
	"os"
)

var RequestLog *log.Logger

func InitLogs() {
	logFile, err := os.OpenFile("./logs/requests.txt", os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	RequestLog = log.New(logFile, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
}
