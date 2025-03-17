package logger

import (
	"log"
	"os"
)

type Logger struct {
  info *log.Logger
  debug *log.Logger
  warn *log.Logger
  err *log.Logger
}

const (
  Reset  = "\033[0m"
  Red    = "\033[31m"
  Green  = "\033[32m"
  Yellow = "\033[33m"
  Blue   = "\033[34m"
)

func New() *Logger {
  return &Logger{}
}

func (l *Logger) Configure() {
  l.info = log.New(os.Stdout, Green+"[INFO] "+Reset, log.Ldate | log.Ltime)
  l.debug = log.New(os.Stdout, Blue+"[DEBUG] "+Reset, log.Ldate | log.Ltime)
  l.warn = log.New(os.Stdout, Yellow+"[WARN] "+Reset, log.Ldate | log.Ltime)

  errFile, err := os.OpenFile("errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
  if err != nil {
    panic(err)
  }
  l.err = log.New(errFile, Red+"[ERROR] "+Reset, log.Ldate | log.Ltime)
}

func (l *Logger) Info(args ...any) {
  l.info.Println(args...)
}

func (l *Logger) Debug(args ...any) {
  l.debug.Println(args...)
}

func (l *Logger) Warn(args ...any) {
  l.warn.Println(args...)
}

func (l *Logger) Error(args ...any) {
  l.err.Println(args...)
}