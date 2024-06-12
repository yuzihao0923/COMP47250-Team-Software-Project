package log

import (
	"log"
)

type logMessage struct {
	level   string
	message string
}

var logChannel = make(chan logMessage, 100)

func init() {
	go processLogMessages()
}

func processLogMessages() {
	for logMsg := range logChannel {
		log.Printf("[%s] %s\n", logMsg.level, logMsg.message)
	}
}

func LogMessage(level string, message string) {
	logChannel <- logMessage{level: level, message: message}
}
