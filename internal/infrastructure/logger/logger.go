package logger

import (
	"context"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	global       *zap.SugaredLogger
	defaultLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
)

func init() {
	SetLogger(NewStdOut(defaultLevel, zap.AddCaller()))
}

func New(level zapcore.LevelEnabler, w io.Writer, options ...zap.Option) *zap.SugaredLogger {
	if level == nil {
		level = defaultLevel
	}

	cfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(zapcore.NewJSONEncoder(cfg), zapcore.AddSync(w), level)
	base := zap.New(core, append(options, zap.AddCaller())...)
	return base.Sugar()
}

func NewStdOut(level zapcore.LevelEnabler, options ...zap.Option) *zap.SugaredLogger {
	return New(level, os.Stdout, options...)
}

func SetLogger(l *zap.SugaredLogger) {
	global = l
}

type ctxKey struct{}

func WithContext(ctx context.Context, keyValues ...any) context.Context {
	log := FromContext(ctx).With(keyValues...)
	return context.WithValue(ctx, ctxKey{}, log)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return global
	}
	if l, ok := ctx.Value(ctxKey{}).(*zap.SugaredLogger); ok && l != nil {
		return l
	}
	return global
}

func Info(msg string, args ...any)  { global.Infow(msg, args...) }
func Error(msg string, args ...any) { global.Errorw(msg, args...) }
func Warn(msg string, args ...any)  { global.Warnw(msg, args...) }
func Debug(msg string, args ...any) { global.Debugw(msg, args...) }
func Fatal(msg string, args ...any) { global.Fatalw(msg, args...) }

func Sync() error {
	if global != nil {
		return global.Sync()
	}
	return nil
}
