package logger

import (
	"context"
	"fmt"
	"github.com/POSIdev-community/aictl/internal/core/domain/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Logger struct {
	z *zap.Logger
}

func NewLogger(verbose bool) (*zap.Logger, error) {
	cores := make([]zapcore.Core, 0, 2)
	cores = append(cores, newInfoCore())

	if verbose {
		cores = append(cores, newErrorCore())
	}

	core := zapcore.NewTee(cores...)

	return zap.New(core, zap.AddCaller()), nil
}

func newInfoCore() zapcore.Core {
	infoEncoderConfig := zapcore.EncoderConfig{
		MessageKey: "msg",
		LineEnding: zapcore.DefaultLineEnding,
	}

	infoEncoder := zapcore.NewConsoleEncoder(infoEncoderConfig)

	return zapcore.NewCore(
		infoEncoder,
		zapcore.Lock(os.Stdout),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
		}),
	)
}

func newErrorCore() zapcore.Core {
	errorEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "",
		CallerKey:      "",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	errorEncoder := zapcore.NewConsoleEncoder(errorEncoderConfig)

	return zapcore.NewCore(
		errorEncoder,
		zapcore.Lock(os.Stderr),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		}),
	)
}

func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(zap.Logger{}).(*zap.Logger); ok {
		return &Logger{logger}
	}

	return &Logger{zap.L()}
}

func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, zap.Logger{}, logger)
}

func (log *Logger) LogConfig(cfg *config.Config) {
	log.z.Error("config",
		zap.String("project-id", cfg.ProjectId().String()),
		zap.String("branch-id", cfg.BranchId().String()))
}

func (log *Logger) StdOut(msg string) {
	log.z.Info(msg)
}

func (log *Logger) StdOutF(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	log.z.Info(msg)
}

func (log *Logger) StdErr(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	log.z.Error(msg)
}
