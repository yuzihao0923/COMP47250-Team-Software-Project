package log

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type logMessage struct {
	level   string
	source  string
	message string
}

var logChannel = make(chan logMessage, 100)
var logEntries []string
var logMutex sync.Mutex

// BroadcastFunc is a function that will be called to broadcast log messages.
var BroadcastFunc func(string)

func init() {
	go processLogMessages()
}

func processLogMessages() {
	for logMsg := range logChannel {
		entry := fmt.Sprintf("[%s] [%s] %s", logMsg.level, logMsg.source, logMsg.message)
		logMutex.Lock()
		logEntries = append(logEntries, entry)
		logMutex.Unlock()
		log.Print(entry)

		// Broadcast the log entry if BroadcastFunc is set
		if BroadcastFunc != nil {
			BroadcastFunc(entry)
		}
	}
}

func LogMessage(level, source, message string) {
	logChannel <- logMessage{level: level, source: source, message: message}
}

func LogInfo(source, message string) {
	LogMessage("INFO", source, message)
}

func LogWarning(source, message string) {
	LogMessage("WARNING", source, message)
}

func LogError(source, message string) {
	LogMessage("ERROR", source, message)
}

func GetLogEntries() []string {
	logMutex.Lock()
	defer logMutex.Unlock()
	return logEntries
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	LogError("HTTP", err.Error())
	http.Error(w, err.Error(), statusCode)
}
