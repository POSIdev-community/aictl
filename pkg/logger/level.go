package logger

import "go.uber.org/zap/zapcore"

//go:generate enumer -type Level -text -transform=lower

type Level uint8

const (
	Debug Level = iota
	Info
	Warn
	Error
	Panic
	Fatal
	DPanic
)

func (l Level) ToZapLevel() zapcore.Level {
	switch l {
	case Debug:
		return zapcore.DebugLevel
	case Info:
		return zapcore.InfoLevel
	case Warn:
		return zapcore.WarnLevel
	case Error:
		return zapcore.ErrorLevel
	case Panic:
		return zapcore.PanicLevel
	case Fatal:
		return zapcore.FatalLevel
	case DPanic:
		return zapcore.DPanicLevel
	default:
		return zapcore.InfoLevel
	}
}
