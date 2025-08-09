package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Logger struct {
	level LogLevel
}

var defaultLogger = &Logger{level: DEBUG}

func init() {
	logLevelStr := os.Getenv("LOG_LEVEL")
	if logLevelStr != "" {
		switch strings.ToUpper(logLevelStr) {
		case "DEBUG":
			defaultLogger.level = DEBUG
		case "INFO":
			defaultLogger.level = INFO
		case "WARN":
			defaultLogger.level = WARN
		case "ERROR":
			defaultLogger.level = ERROR
		case "FATAL":
			defaultLogger.level = FATAL
		default:
			log.Printf("Invalid LOG_LEVEL environment variable: %s, using default INFO level", logLevelStr)
		}
	}
}

func SetLogLevel(level LogLevel) {
	defaultLogger.level = level
}

func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level >= l.level {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		levelStr := l.levelString(level)
		message := fmt.Sprintf(format, args...)
		log.Printf("[%s] %s: %s", timestamp, levelStr, message)
	}
}

func (l *Logger) levelString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
	os.Exit(1)
}

func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	defaultLogger.Fatal(format, args...)
}

func LogError(message string, err error) {
	if err != nil {
		Error("%s: %v", message, err)
	} else {
		Error(message)
	}
}
