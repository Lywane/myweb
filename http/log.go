package http

import (
	syslog "log"
	"os"
	"fmt"
)

type loggerer interface {
	Log(tag string, msg ...interface{})
}

type cmdLogger struct{}

type fileLogger struct {
	Dir  string `json:"dir"`
	Name string `json:"name"`
}

var log loggerer

func init() {
	log = &cmdLogger{}
}

func SetLogger(logger loggerer) {
	log = logger
}

func (clg *cmdLogger) Log(tag string, msg ...interface{}) {
	syslog.Println(append([]interface{}{tag}, msg...)...)
}

func NewFileLogger(name, dir string) *fileLogger {
	return &fileLogger{
		Dir:  dir,
		Name: name,
	}
}

func (flg *fileLogger) Log(tag string, msg ...interface{}) {
	f, err := os.OpenFile(flg.Dir+flg.Name, os.O_CREATE|os.O_APPEND, 0x666)
	if err != nil {
		panic(err)
	}
	loger := syslog.New(f, "", syslog.Ldate|syslog.Ltime)
	message := tag
	for _, m := range msg {
		message += fmt.Sprintf(" %v", m)
	}
	loger.Println(message)
}

func Info(msg ...interface{}) {
	log.Log("INFO", msg...)
}

func Warn(msg ...interface{}) {
	log.Log("WARN", msg...)
}

func Error(msg ...interface{}) {
	log.Log("ERROR", msg...)
}

func Debug(msg ...interface{}) {
	log.Log("DEBUG", msg...)
}
