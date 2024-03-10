package log

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"time"
)

func newZapLog() Logger {
	dir, _ := os.Getwd()
	fmt.Printf("NewZapLog dir:%s \n", dir)
	encoder := zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
				encoder.AppendString(t.Format("2006-01-02 15:04:05"))
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
