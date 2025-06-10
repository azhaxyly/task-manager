package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var logInstance *log.Logger

func Init(out io.Writer) {
	logInstance = log.New(out, "", log.LstdFlags|log.Lmicroseconds)
}

func Info(format string, args ...interface{}) {
	if logInstance == nil {
		logInstance = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	}
	prefix := fmt.Sprintf("%s INFO: ", time.Now().Format("2006/01/02 15:04:05.000000"))
	logInstance.Output(2, prefix+fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	if logInstance == nil {
		logInstance = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	}
	prefix := fmt.Sprintf("%s ERROR: ", time.Now().Format("2006/01/02 15:04:05.000000"))
	logInstance.Output(2, prefix+fmt.Sprintf(format, args...))
}
