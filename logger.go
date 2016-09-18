package JLogger

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"time"
)

const (
	flag_info  = "[info ]"
	flag_debug = "[debug]"
	flag_err   = "[error]"
)

type Logger struct {
	depth     int
	path      string
	pLogger   *log.Logger
	svrName   string
	funcQueue chan func()
}

func New(path string, svrName string, depth int) *Logger {
	pLog := &Logger{
		depth:     depth,
		path:      path,
		svrName:   svrName,
		pLogger:   log.New(nil, "", log.LstdFlags),
		funcQueue: make(chan func(), 128),
	}
	pLog.dayTimer(pLog.refreshWriter)
	go func() {
		for {
			f := <-pLog.funcQueue
			f()
		}
	}()
	return pLog
}

//每日循环
func (l *Logger) dayTimer(f func()) {
	f()
	now := time.Now().Local()
	y, m, d := now.Date()
	nextDate := time.Date(y, m, d+1, 0, 0, 0, 0, time.Local)
	delayDuration := nextDate.Sub(now)
	go func() {
		<-time.After(delayDuration)
		l.dayTimer(f)
	}()
}

//刷新writer
func (l *Logger) refreshWriter() {
	fileName := l.getFileName()
	strPath := path.Join(l.path, fileName)
	pFile, err := os.OpenFile(strPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		print(strPath + " createFile err:" + err.Error())
		return
	}
	l.pLogger.SetOutput(pFile)
}

func (l *Logger) getFileName() string {
	y, m, d := time.Now().Date()
	fileName := fmt.Sprintf("%v_%d-%02d-%02d.log", l.svrName, y, m, d)
	return fileName
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.print(flag_info, format, args...)
}
func (l *Logger) Debug(format string, args ...interface{}) {
	l.print(flag_debug, format, args...)
}
func (l *Logger) Err(format string, args ...interface{}) {
	l.print(flag_err, format, args...)
}

func (l *Logger) print(flag string, format string, args ...interface{}) {
	stack := "??"
	if _, file, line, ok := runtime.Caller(l.depth); ok {
		stack = fmt.Sprintf("[%v:%v]", file, line)
	}
	format = stack + flag + format
	l.funcQueue <- func() { l.pLogger.Printf(format, args...) }
}
