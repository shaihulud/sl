package sl

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"log"
	"os"
)

// Свой тип логгера, который не переоткрывает файлы и имеет удобный формат
type SmartLogger struct {
	logger  *log.Logger
	fileObj *os.File
}

// Создает глобальный логгер
func CreateSmartLogger(logsDir, logName string) *SmartLogger {
	filePath := fmt.Sprintf("%s/%s", logsDir, logName)
	fileObj, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	logger := log.New(fileObj, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	return &SmartLogger{logger, fileObj}
}

// Добавляет в лог запись с шапкой INFO
func (l *SmartLogger) Info(format string, v ...interface{}) {
	l.logger.Printf("INFO \t" + format + "\n", v...)
}

// Добавляет в лог запись с шапкой WARN
func (l *SmartLogger) Warning(format string, tags map[string]string, v ...interface{}) {
	l.logger.Printf("WARNING\t" + format + "   %v\n", tags)
	raven.CaptureMessage(format, tags)
}

// Добавляет в лог запись с шапкой ERROR
func (l *SmartLogger) Error(format string, v ...interface{}) {
	for _, i := range v {
		switch obj := i.(type) {
		case error:
			raven.CaptureError(obj.(error), nil)
		}
	}
	l.logger.Printf("ERROR\t" + format + "\n", v...)
}

// Добавляет в лог запись с шапкой FATAL и выходит
func (l *SmartLogger) Fatal(format string, v ...interface{}) {
	for _, i := range v {
		switch obj := i.(type) {
		case error:
			raven.CaptureError(obj.(error), nil)
		}
	}
	l.logger.Fatalf("FATAL\t" + format + "\n", v...)
}

// Закрывает глобальный логгер
func (l *SmartLogger) Close() error {
	l.fileObj.Sync()
	return l.fileObj.Close()
}
