package logger_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	. "github.com/POSIdev-community/aictl/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newFixedClock(t time.Time) *fixedClock {
	return &fixedClock{t: t}
}

type fixedClock struct {
	t time.Time
}

func (f *fixedClock) Now() time.Time {
	return f.t
}

func (f *fixedClock) NewTicker(duration time.Duration) *time.Ticker {
	return time.NewTicker(duration)
}

func NewCustomLogger(pipeTo io.Writer) zapcore.Core {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.ConsoleSeparator = "  "

	enc, _ := NewVerboseEncoder(cfg.EncoderConfig)

	return zapcore.NewCore(
		enc,
		zapcore.AddSync(pipeTo),
		zapcore.InfoLevel,
	)
}

func TestInitDevLogger(t *testing.T) {
	t.Parallel()

	log := InitDevLogger()
	buf := bytes.NewBuffer(nil)

	// Fixate date
	clock := newFixedClock(time.Date(2025, 1, 16, 23, 15, 0, 0, time.UTC))
	log = log.WithOptions(zap.WithClock(clock))
	// Write to buf instead of stderr
	log = log.WithOptions(zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return NewCustomLogger(buf)
	}))

	log = log.With(zap.String("foo", "bar"))
	log.Info("hello")
	log2 := log.With(zap.String("baz", "bar"))
	log2.Info("world", zap.Int("int", 123))
	log.Info("helloworld")

	want := `2025-01-16T23:15:00.000Z  INFO   hello
  foo = bar
2025-01-16T23:15:00.000Z  INFO   world
  baz = bar
  foo = bar
  int = 123
2025-01-16T23:15:00.000Z  INFO   helloworld
  foo = bar
`
	assert.Equal(t, want, buf.String())
}
