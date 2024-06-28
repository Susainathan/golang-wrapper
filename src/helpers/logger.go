package helpers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	logFile *os.File
	logger  *log.Logger
	logLock sync.Mutex
)

func InitLogger() {
	var err error
	logFile, err = os.OpenFile(filepath.Base(".")+"/logs/logs.logs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// logFile, err = os.OpenFile("/tnq/apps/xml-central/xmlc-wrapper-go/logs/logs.logs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal("Error opening log file: ", err)
	}
	logger = log.New(logFile, "", log.LstdFlags)
}

func LogMessage(level string, message string) {
	logLock.Lock()
	defer logLock.Unlock()
	lMessage := fmt.Sprintf("[%s] %s\n", level, message)
	logger.Printf(lMessage)
	log.Println(lMessage)

}
