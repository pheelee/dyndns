package logger

import (
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime)
}

// Info emits info severity log message
func Info(msg string) {
	log.Println("[INFO] " + msg)
}

// Warning emits warning severity log message
func Warning(msg string) {
	log.Println("[WARNING] " + msg)
}

// Error emits error severity log message
func Error(msg string) {
	log.Println("[ERROR] " + msg)
}
