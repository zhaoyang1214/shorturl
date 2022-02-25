package logger

import (
	"github.com/zhaoyang1214/ginco/framework/contract"
	"go.uber.org/zap"
	"strconv"
)

type Logger struct {
	*zap.Logger
}

var _ contract.Logger = (*Logger)(nil)

func (l *Logger) convert(context ...interface{}) []zap.Field {
	fields := make([]zap.Field, len(context))
	for i, field := range context {
		switch field.(type) {
		case zap.Field:
			fields[i] = field.(zap.Field)
		default:
			fields[i] = zap.Reflect(strconv.Itoa(i), field)
		}
	}
	return fields
}

func (l *Logger) Debug(msg string, context ...interface{}) {
	l.Logger.Debug(msg, l.convert(context...)...)
}

func (l *Logger) Info(msg string, context ...interface{}) {
	l.Logger.Info(msg, l.convert(context...)...)
}

func (l *Logger) Warn(msg string, context ...interface{}) {
	l.Logger.Warn(msg, l.convert(context...)...)
}

func (l *Logger) Error(msg string, context ...interface{}) {
	l.Logger.Error(msg, l.convert(context...)...)
}

func (l *Logger) Panic(msg string, context ...interface{}) {
	l.Logger.Panic(msg, l.convert(context...)...)
}

func (l *Logger) Fatal(msg string, context ...interface{}) {
	l.Logger.Fatal(msg, l.convert(context...)...)
}

func (l *Logger) Log(level interface{}, msg string, context ...interface{}) {
	if ce := l.Logger.Check(level.(Level), msg); ce != nil {
		ce.Write(l.convert(context...)...)
	}
}

func (l *Logger) Sync() error {
	return l.Logger.Sync()
}
