package types

import "log"

type LogSvc interface {
	GetErr() *log.Logger
	GetOut() *log.Logger
}
