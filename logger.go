package gorelay

import "log"

type Logger interface {
	Info(message string)
	Error(message string)
}

type BasicLogger struct {
}

func (l *BasicLogger) Info(message string) {
	log.Println("INFO: " + message)
}

func (l *BasicLogger) Error(message string) {
	log.Println("ERROR: " + message)
}

func NewBasicLogger() Logger {
	return &BasicLogger{}
}
