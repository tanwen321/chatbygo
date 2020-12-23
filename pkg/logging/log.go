package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Level int

var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func init() {
	filePath := getLogFileFullPath()
	F = openLogFile(filePath)

	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

// func Printf(format string, v ...interface{}) {
// 	setPrefix(INFO)
// 	logger.Printf(format, v)
// }

func Debug(f string, v ...interface{}) {
	setPrefix(DEBUG)
	logger.Printf(f, v...)
}

func Info(f string, v ...interface{}) {
	setPrefix(INFO)
	logger.Printf(f, v...)
}

func Warn(f string, v ...interface{}) {
	setPrefix(WARNING)
	logger.Printf(f, v...)
}

func Error(f string, v ...interface{}) {
	setPrefix(ERROR)
	logger.Printf(f, v...)
}

func Fatal(f string, v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalf(f, v...)
}

func Fatalf(f string, v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalf(f, v...)
}

func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}
