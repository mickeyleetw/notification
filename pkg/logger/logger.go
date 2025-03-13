package logger

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "NOTIFY-SYSTEM: ", log.LstdFlags)

// Info is a function that logs an info message
func Info(msg string) {
	logger.Println("INFO: " + msg)
}

// Error is a function that logs an error message
func Error(msg string) {
	logger.Println("ERROR: " + msg)
}

// Fatal is a function that logs a fatal message and then terminates the program
func Fatal(msg string) {
	logger.Fatal("FATAL: " + msg)
}
