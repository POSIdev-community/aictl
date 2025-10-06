package logger

import (
	"context"
	"fmt"
	"log"
	"maps"
	"slices"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

const (
	defaultProdLogLevel = "info"
	defaultDevLogLevel  = "debug"

	timeFormat = "2006-01-02T15:04:05.000Z07:00"
)

type Config struct {
	Level string `env:"LEVEL, default=info"`
}

func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(zap.Logger{}).(*zap.Logger); ok {
		return logger
	}

	return zap.L()
}

func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, zap.Logger{}, logger)
}

func InitProdLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	zapLevel := getZapLevel(defaultProdLogLevel)
	cfg.Level = zap.NewAtomicLevelAt(zapLevel)

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("create logger: %v", err)
	}

	zap.ReplaceGlobals(logger)
	zap.RedirectStdLog(logger)

	return logger
}

func InitDevLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()

	if err := zap.RegisterEncoder("verbose", NewVerboseEncoder); err != nil {
		log.Fatalf("Register logger encoder failed: %v", err)
	}

	cfg.Encoding = "verbose"
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true
	cfg.EncoderConfig.ConsoleSeparator = "  "

	zapLevel := getZapLevel(defaultDevLogLevel)
	cfg.Level = zap.NewAtomicLevelAt(zapLevel)

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("create logger: %v", err)
	}

	zap.ReplaceGlobals(logger)
	zap.RedirectStdLog(logger)

	return logger
}

type verboseEncoder struct {
	*zapcore.MapObjectEncoder
	cfg zapcore.EncoderConfig
}

// EncodeEntry преобразует запись лога и поля в кастомный формат.
func (e *verboseEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// Создаем MapObjectEncoder для кодирования полей
	enc := zapcore.NewMapObjectEncoder()
	for _, field := range fields {
		field.AddTo(enc)
	}

	for name, value := range e.MapObjectEncoder.Fields {
		enc.Fields[name] = value
	}

	// Формируем строку в кастомном формате
	var builder strings.Builder

	// Добавляем уровень и сообщение
	builder.WriteString(entry.Time.Format(timeFormat))
	builder.WriteString(e.cfg.ConsoleSeparator)
	builder.WriteString(entry.Level.CapitalString())
	builder.WriteString(e.cfg.ConsoleSeparator)

	if entry.Level != zapcore.DebugLevel && entry.Level != zapcore.ErrorLevel {
		builder.WriteString(" ")
	}

	if entry.LoggerName != "" {
		builder.WriteString(entry.LoggerName)
		builder.WriteString(e.cfg.ConsoleSeparator)
	}

	if entry.Caller.Defined {
		builder.WriteString(entry.Caller.TrimmedPath())
		builder.WriteString(e.cfg.ConsoleSeparator)
	}

	builder.WriteString(entry.Message)

	var fieldNameMaxLength int

	fieldNames := slices.AppendSeq(make([]string, 0, len(e.Fields)), maps.Keys(enc.Fields))

	slices.Sort(fieldNames)

	for _, name := range fieldNames {
		fieldNameMaxLength = max(fieldNameMaxLength, len(name))
	}

	// Добавляем поля в формате key=value
	for _, key := range fieldNames {
		builder.WriteString(e.cfg.LineEnding)
		builder.WriteString(e.cfg.ConsoleSeparator)
		builder.WriteString(key)
		builder.WriteString(zFill(fieldNameMaxLength - len(key)))
		builder.WriteString(" = ")
		builder.WriteString(fmt.Sprintf("%v", enc.Fields[key]))
	}

	// Добавляем перевод строки
	builder.WriteString(e.cfg.LineEnding)

	// Возвращаем буфер с данными
	buf := buffer.NewPool().Get()
	buf.AppendString(builder.String())

	return buf, nil
}

func (e *verboseEncoder) Clone() zapcore.Encoder {
	enc := zapcore.NewMapObjectEncoder()
	maps.Copy(enc.Fields, e.Fields)

	return &verboseEncoder{
		MapObjectEncoder: enc,
		cfg:              e.cfg,
	}
}

func NewVerboseEncoder(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
	return &verboseEncoder{
		MapObjectEncoder: zapcore.NewMapObjectEncoder(),
		cfg:              cfg,
	}, nil
}

func zFill(length int) string {
	if length == 0 {
		return ""
	}

	s := make([]rune, length)
	for i := range s {
		s[i] = ' '
	}

	return string(s)
}

func getZapLevel(level string) zapcore.Level {
	lvl, err := LevelString(level)
	if err != nil {
		return zapcore.InfoLevel
	}

	return lvl.ToZapLevel()
}
