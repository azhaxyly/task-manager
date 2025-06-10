package logger

import (
	"io"
	"log"
	"os"
)

var logInstance *log.Logger

func Init(out io.Writer) {
	logInstance = log.New(out, "", log.LstdFlags|log.Lmicroseconds)
}

func Info(format string, args ...interface{}) {
	if logInstance == nil {
		logInstance = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	}
	logInstance.Printf("INFO: "+format, args...)
}

func Error(format string, args ...interface{}) {
	if logInstance == nil {
		logInstance = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	}
	logInstance.Printf("ERROR: "+format, args...)
}
