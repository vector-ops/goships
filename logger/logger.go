package logger

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/vector-ops/goships/types"
)

type LoggerOption func(logger *Logger) error

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

// Destinations
const (
	File = iota
	Channel
)

type Logger struct {
	logCh chan Log

	file *os.File

	dst []int
}

type Log struct {
	Timestamp time.Time `json:"timestamp"`
	Msg       string    `json:"msg"`
	Level     Level     `json:"level"`
	Color     int16     `json:"color"`
}

func NewLogger(opts ...LoggerOption) (*Logger, error) {

	if len(opts) == 0 {
		return nil, errors.New("no options provided")
	}

	logger := new(Logger)

	logger.dst = make([]int, 0)

	for _, opt := range opts {
		if err := opt(logger); err != nil {
			return nil, err
		}
	}

	return logger, nil

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

	l.writeToAllDst(log)

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

	l.writeToAllDst(log)
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

	l.writeToAllDst(log)
}

func (l *Logger) writeToAllDst(log Log) error {

	for _, dst := range l.dst {
		switch dst {
		case Channel:
			l.logCh <- log
		case File:
		}
	}

	return nil
}

func WithLogChan(logCh chan Log) LoggerOption {
	return func(logger *Logger) error {
		if logCh == nil {
			return errors.New("log channel cannot be nil")
		}

		logger.logCh = logCh

		logger.dst = append(logger.dst, Channel)

		return nil
	}
}

func WithLogFile(fp string) LoggerOption {
	return func(logger *Logger) error {

		file, err := os.OpenFile(fp, os.O_RDWR, os.ModePerm)
		if err != nil {
			return err
		}

		logger.file = file

		logger.dst = append(logger.dst, File)

		return nil
	}
}

func WithNoOp() LoggerOption {
	return func(logger *Logger) error {
		return nil
	}
}
