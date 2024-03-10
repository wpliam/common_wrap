package log

func Debug(args ...interface{}) {
	GetDefaultLogger().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	GetDefaultLogger().Debugf(format, args...)
}

func Info(args ...interface{}) {
	GetDefaultLogger().Info(args...)
}

func Infof(format string, args ...interface{}) {
	GetDefaultLogger().Infof(format, args...)
}

func Error(args ...interface{}) {
	GetDefaultLogger().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	GetDefaultLogger().Errorf(format, args...)
}
