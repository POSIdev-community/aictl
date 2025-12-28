package logger

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/POSIdev-community/aictl/internal/core/domain/config"
)

type Logger struct {
	z *zap.Logger
}

func NewLogger(consoleVerbose bool, logPath string) (*zap.Logger, error) {
	cores := make([]zapcore.Core, 0, 3)
	cores = append(cores, newInfoCore())

	if consoleVerbose {
		cores = append(cores, newErrorCore())
	}

	if logPath != "" {
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file %q: %w", logPath, err)
		}

		fileEncoder := zapcore.NewConsoleEncoder(newErrorConfig())

		fileCore := zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(file),
			zapcore.DebugLevel,
		)

		cores = append(cores, fileCore)
	}

	core := zapcore.NewTee(cores...)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)), nil
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
	errorEncoder := zapcore.NewConsoleEncoder(newErrorConfig())

	return zapcore.NewCore(
		errorEncoder,
		zapcore.Lock(os.Stderr),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		}),
	)
}

func newErrorConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
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
	log.z.Sugar().Info(msg)
}

func (log *Logger) StdOutf(format string, a ...any) {
	log.z.Sugar().Infof(format, a...)
}

func (log *Logger) StdErr(msg string) {
	log.z.Sugar().Error(msg)
}

func (log *Logger) StdErrf(format string, a ...any) {
	log.z.Sugar().Errorf(format, a...)
}
