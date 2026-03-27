package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/vector-ops/goships/types"
)

type Level string

func (l Level) String() string {
	return string(l)
}

func (l Level) Upper() string {
	return strings.ToUpper(l.String())
}

const (
	Info  Level = "info"
	Warn  Level = "warn"
	Error Level = "error"
)

type Logger struct {
	logCh chan Log

	writeToFile bool
	file        *os.File
}

type Log struct {
	Timestamp time.Time `json:"timestamp"`
	Msg       string    `json:"msg"`
	Level     Level     `json:"level"`
	Color     int16     `json:"color"`
}

func NewLogger(logCh chan Log, writeToFile bool, fp *string) *Logger {

	var file *os.File
	var err error

	if writeToFile && fp != nil {
		file, err = os.OpenFile(*fp, os.O_RDWR, os.ModePerm)
		if err != nil {
			panic(err)
		}
	} else {
		writeToFile = false
	}

	return &Logger{
		logCh: logCh,

		writeToFile: writeToFile,
		file:        file,
	}
}

func (l *Logger) Infof(msg string, vars ...any) {

	if len(vars) > 0 {
		msg = fmt.Sprintf(msg, vars)
	}

	if len(msg) == 0 {
		return
	}

	log := Log{
		Timestamp: time.Now(),
		Msg:       msg,
		Level:     Info,
		Color:     types.WhiteBlack,
	}

	l.logCh <- log

}

func (l *Logger) Errorf(msg string, vars ...any) {

	if len(vars) > 0 {
		msg = fmt.Sprintf(msg, vars)
	}

	if len(msg) == 0 {
		return
	}

	log := Log{
		Timestamp: time.Now(),
		Msg:       msg,
		Level:     Error,
		Color:     types.RedBlack,
	}

	l.logCh <- log
}

func (l *Logger) Warnf(msg string, vars ...any) {

	if len(vars) > 0 {
		msg = fmt.Sprintf(msg, vars)
	}

	if len(msg) == 0 {
		return
	}

	log := Log{
		Timestamp: time.Now(),
		Msg:       msg,
		Level:     Warn,
		Color:     types.YellowBlack,
	}

	l.logCh <- log
}
