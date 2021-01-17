package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var MainLogger = newLogger()

type logger struct {
	sync.Mutex
	logFile *os.File
}

func newLogger() *logger {
	log := &logger{}
	log.createFile()

	return log
}

func (l *logger) createFile() {
	fileName := time.Now().Format("2006-01-02T15-04-05")
	fileName = strings.Replace(fileName, "T", "-", 1)

	absPath := "./log"
	err := os.MkdirAll(absPath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	l.logFile, err = os.OpenFile(
		filepath.Join(absPath, "log-"+fileName+".log"),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0600)
	if err != nil {
		fmt.Println(err)
	}
}

func (l *logger) Log(data string) {
	l.Lock()
	data = l.prepareData(data)
	fmt.Println(data)
	l.saveToFile(data)
	l.Unlock()
}

func (l *logger) prepareData(data string) string {
	timeStamp := time.Now().Format("2006-01-02 15:04:05.000")

	return fmt.Sprintf("[%s] %s", timeStamp, data)
}

func (l *logger) saveToFile(data string) {
	if l.logFile != nil {
		_, _ = l.logFile.WriteString(data + "\n")
	}
}

func (l *logger) CloseLogFile() {
	l.Lock()
	if l.logFile != nil {
		_ = l.logFile.Close()
	}
	l.Unlock()
}
