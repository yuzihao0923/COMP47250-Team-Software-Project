package log

import (
	"log"
)

func LogMessage(level string, message string) {
	// currentTime := time.Now().Format("2006/01/02 15:04:05")
	// timestamp := time.Now().Format(time.RFC3339)
	log.Printf("[%s] %s\n", level, message)
}
