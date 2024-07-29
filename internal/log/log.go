package log

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type logMessage struct {
	level   string
	source  string
	message string
}

var (
	logChannel    = make(chan logMessage, 100)
	logEntries    []string
	logMutex      sync.Mutex
	logFile       *os.File
	BroadcastFunc func(string) // BroadcastFunc is a function that will be called to broadcast log messages.
)

func init() {
	// 确保 logs 目录存在
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", os.ModePerm)
	}

	// 打开日志文件
	var err error
	logFile, err = os.OpenFile("/home/yuzihao0923/COMP47250-Team-Software-Project/tests/Redis-throuput/logs/broker.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	go processLogMessages()
}

func processLogMessages() {
	for logMsg := range logChannel {
		entry := fmt.Sprintf("[%s] [%s] %s", logMsg.level, logMsg.source, logMsg.message)
		logMutex.Lock()
		logEntries = append(logEntries, entry)
		logMutex.Unlock()
		log.Print(entry)

		// 写入文件
		if _, err := logFile.WriteString(entry + "\n"); err != nil {
			log.Fatal(err)
		}

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
