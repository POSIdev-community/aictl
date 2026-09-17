package logger

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Console verbosity levels for NewLogger.
const (
	LevelQuiet   = 0 // Info (stdout) only
	LevelVerbose = 1 // + Error ops on stderr
	LevelDebug   = 2 // + Debug (error chains) on stderr
)

type contextKey struct{}

type Logger struct {
	z     *zap.Logger
	file  *zap.Logger // optional file-only sink (Error / Debug); nil if no --log-path
	level int
}

// Wrap builds a Logger around an existing zap logger (tests / special sinks).
func Wrap(z *zap.Logger) *Logger {
	return &Logger{z: z, level: LevelQuiet}
}

func NewLogger(verboseLevel int, logPath string) (*Logger, error) {
	cores := make([]zapcore.Core, 0, 4)
	cores = append(cores, newInfoCore())

	if verboseLevel >= LevelVerbose {
		cores = append(cores, newErrorCore())
	}
	if verboseLevel >= LevelDebug {
		cores = append(cores, newDebugCore())
	}

	var fileLogger *zap.Logger
	if logPath != "" {
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file %q: %w", logPath, err)
		}

		fileEncoder := zapcore.NewConsoleEncoder(newErrorConfig())
		fileCore := zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(file),
			fileLevelEnabler(verboseLevel >= LevelDebug),
		)

		cores = append(cores, fileCore)
		fileLogger = zap.New(fileCore, zap.AddCaller(), zap.AddCallerSkip(1))
	}

	core := zapcore.NewTee(cores...)
	z := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &Logger{z: z, file: fileLogger, level: verboseLevel}, nil
}

// IsVerbose reports whether -v/--verbose or -V/--debug is enabled.
func (log *Logger) IsVerbose() bool {
	return log != nil && log.level >= LevelVerbose
}

func fileLevelEnabler(includeDebug bool) zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		if lvl >= zapcore.ErrorLevel {
			return true
		}
		return includeDebug && lvl == zapcore.DebugLevel
	}
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

func newDebugCore() zapcore.Core {
	debugEncoder := zapcore.NewConsoleEncoder(newErrorConfig())

	return zapcore.NewCore(
		debugEncoder,
		zapcore.Lock(os.Stderr),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.DebugLevel && lvl < zapcore.ErrorLevel
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
	if l, ok := ctx.Value(contextKey{}).(*Logger); ok {
		return l
	}

	return &Logger{z: zap.L()}
}

func ContextWithLogger(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, l)
}

func (log *Logger) LogConfig(projectID, branchID string) {
	log.z.Error("config",
		zap.String("project-id", projectID),
		zap.String("branch-id", branchID))
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

func (log *Logger) Debugf(format string, a ...any) {
	log.z.Sugar().Debugf(format, a...)
}

// FileError writes msg to the log file only (zap Error + timestamp). No-op without --log-path.
// Does not write to the console, so it is safe alongside fmt.Fprintln to stderr.
func (log *Logger) FileError(msg string) {
	if log.file == nil {
		return
	}
	log.file.Sugar().Error(msg)
}

func (log *Logger) HasFile() bool {
	return log.file != nil
}
