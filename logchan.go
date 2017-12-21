package logger

import (
	"fmt"
	"log"
	"os"
	"time"
	"io"
)

const (
	_      = iota
	DEBUG
	INFO
	NOTICE
	WARN
	ERROR
	FATAL
)

var lineHead = map[int]string{
	DEBUG:  "D",
	INFO:   "I",
	NOTICE: "N",
	WARN:   "W",
	ERROR:  "E",
	FATAL:  "F",
}

type LogType struct {
	level int
	log   string
}

var logChannel chan LogType
var logPath, logPrefix string

func StartLog(path, prefix string, buflen int, done chan struct{}) {
	logChannel = make(chan LogType, buflen)
	go logging(path, prefix, done)
}

func Log(lv int, format string, a ...interface{}) {
	logChannel <- LogType{level: lv, log: fmt.Sprintf(format, a...)}
}

func Close() {
	time.Sleep(time.Millisecond * 100)
	close(logChannel)
}

func rotateLogFile(t time.Time, f *os.File) *os.File {
	if f != nil {
		f.Close()
	}

	year, mon, day := t.Date()
	fpLog, err := os.OpenFile(logPath+"/"+logPrefix+fmt.Sprintf("_%04d%02d%02d.log", year, mon, day)+"", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	multiWriter := io.MultiWriter(fpLog, os.Stdout)
	log.SetOutput(multiWriter)
	//log.SetOutput(fpLog)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	return fpLog
}

func logging(path, prefix string, done chan struct{}) {

	now := time.Now()
	nowDay := now.YearDay()

	logPath = path
	logPrefix = prefix

	fpLog := rotateLogFile(now, nil)
	defer func(f *os.File) {
		fpLog.Close()
		done <- struct{}{}
	}(fpLog)

	for log2Write := range logChannel {
		now = time.Now()
		if nowDay != now.YearDay() {
			nowDay = now.YearDay()
			fpLog = rotateLogFile(now, fpLog)
		}

		log.Printf("[%s] %s\n", lineHead[log2Write.level], log2Write.log)
	}

	log.Printf("[%s] %s\n", lineHead[NOTICE], "log channel closed.")
}
