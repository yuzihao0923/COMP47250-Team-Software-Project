package log

import (
    "log"
	"fmt"
    "net/http"
    "sync"
)

type logMessage struct {
    level   string
    message string
}

var logChannel = make(chan logMessage, 100)
var logEntries []string
var logMutex sync.Mutex

func init() {
    go processLogMessages()
}

func processLogMessages() {
    for logMsg := range logChannel {
        entry := fmt.Sprintf("[%s] %s", logMsg.level, logMsg.message)
        logMutex.Lock()
        logEntries = append(logEntries, entry)
        logMutex.Unlock()
        log.Print(entry)
    }
}

func LogMessage(level string, message string) {
    logChannel <- logMessage{level: level, message: message}
}

func LogInfo(message string) {
    LogMessage("INFO", message)
}

func LogError(err error) {
    LogMessage("ERROR", err.Error())
}

func LogWarning(message string) {
    LogMessage("WARNING", message)
}

func GetLogEntries() []string {
    logMutex.Lock()
    defer logMutex.Unlock()
    return logEntries
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, err error) {
    LogError(err)
    http.Error(w, err.Error(), statusCode)
}
