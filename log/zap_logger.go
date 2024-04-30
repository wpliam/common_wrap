package log

import (
	"github.com/natefinch/lumberjack"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newZapLog() Logger {
	// TODO linux不这样取
	dir, _ := os.Getwd()
	encoder := zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeTime: func(t time.Time, encode zapcore.PrimitiveArrayEncoder) {
				encode.AppendString(t.Format("2006-01-02 15:04:05"))
			},
			EncodeLevel: zapcore.CapitalLevelEncoder,
		},
	)
	writeSyncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(dir, "default.log"),
			MaxSize:    500,
			MaxAge:     5,
			MaxBackups: 30,
			LocalTime:  true,
			Compress:   false,
		}),
		zapcore.AddSync(os.Stdout),
	)
	return &zapLog{
		l: zap.New(
			zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel),
			zap.AddCallerSkip(2),
			zap.AddCaller(),
		),
	}
}

type zapLog struct {
	l *zap.Logger
}

func (z *zapLog) Debug(args ...interface{}) {
	z.l.Sugar().Debug(args...)
}

func (z *zapLog) Debugf(format string, args ...interface{}) {
	z.l.Sugar().Debugf(format, args...)
}

func (z *zapLog) Info(args ...interface{}) {
	z.l.Sugar().Info(args...)
}

func (z *zapLog) Infof(format string, args ...interface{}) {
	z.l.Sugar().Infof(format, args...)
}

func (z *zapLog) Error(args ...interface{}) {
	z.l.Sugar().Error(args...)
}

func (z *zapLog) Errorf(format string, args ...interface{}) {
	z.l.Sugar().Errorf(format, args...)
}

func (z *zapLog) Sync() error {
	return z.l.Sugar().Sync()
}
