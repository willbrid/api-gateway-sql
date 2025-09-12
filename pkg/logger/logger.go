package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type LogLevel int

const (
	Info LogLevel = iota
	Warning
	Error
	Fatal
)

const DefaultCallerDepth = 2

var logLevels = map[LogLevel]string{
	Info:    "INFO",
	Warning: "WARNING",
	Error:   "ERROR",
	Fatal:   "FATAL",
}

func msg(level LogLevel, message string, args ...any) {
	logLevel, ok := logLevels[level]
	if !ok {
		logLevel = "UNKNOWN"
	}

	datetime := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	filename := ""
	if ok {
		filename = filepath.Base(file)
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	log.Printf("[%s][%s][%s:%d] : %s", datetime, logLevel, filename, line, fmt.Sprintf(message, args...))
}

func LogInfo(message string, args ...any) {
	msg(Info, message, args...)
}

func LogWarn(message string, args ...any) {
	msg(Warning, message, args...)
}

func LogError(message string, args ...any) {
	msg(Error, message, args...)
}

func LogFatal(message string, args ...any) {
	msg(Fatal, message, args...)

	os.Exit(1)
}
