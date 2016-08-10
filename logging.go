package invt

import (
	"sync"

	"time"

	tl "github.com/apourchet/investment/lib/tradelogger"
)

const (
	TAG_CLOSEPOSITION = "CLOSE_POSITION"
)

var loggers []tl.Logger
var once sync.Once
var lock sync.Mutex

func AddLogger(l tl.Logger) {
	if l == nil {
		return
	}
	once.Do(createLoggers)
	lock.Lock()
	loggers = append(loggers, l)
	lock.Unlock()
}

func Log(date time.Time, tag string, message string) {
	for _, l := range loggers {
		l.Log(&tl.Item{date, tag, message})
	}
}

func createLoggers() {
	lock.Lock()
	loggers = make([]tl.Logger, 0)
	lock.Unlock()
}
