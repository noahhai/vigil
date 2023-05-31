package services

import (
	"github.com/noahhai/vigil/agent/types"
	"log"
	"os"
)

type logger struct {
	out *log.Logger
	err *log.Logger
}

func NewLogger() types.LogSvc {
	return &logger{
		out: log.New(os.Stdout, "", 0),
		err: log.New(os.Stderr, "", 0),
	}
}

func (l *logger) GetOut() *log.Logger {
	return l.out
}
func (l *logger) GetErr() *log.Logger {
	return l.err
}