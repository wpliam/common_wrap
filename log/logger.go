package log

import "sync"

const defaultLogName = "default"

var loggers = make(map[string]Logger)
var mu sync.RWMutex

func init() {
	Register(defaultLogName, newZapLog())
}

func Register(name string, l Logger) {
	mu.Lock()
	defer mu.Unlock()
	_, ok := loggers[name]
	if ok && name == defaultLogName {
		return
	}
	loggers[name] = l
}

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Sync() error
}

func GetDefaultLogger() Logger {
	return Get(defaultLogName)
}

func Get(name string) Logger {
	mu.RLock()
	l := loggers[name]
	mu.RUnlock()
	return l
}

func Sync() {
	mu.RLock()
	defer mu.RUnlock()
	for _, l := range loggers {
		_ = l.Sync()
	}
}
