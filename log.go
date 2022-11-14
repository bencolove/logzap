// refers
// github.com/bigwhite/experiments/tree/master/uber-zap-advanced-usage/demo1/pkg/log/log.go

package logzap

import (
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

const (
	InfoLevel  Level = zap.InfoLevel  // 1
	ErrorLevel Level = zap.ErrorLevel // 2

	// PanicLevel logs message and panics
	PanicLevel Level = zap.PanicLevel //4
	// FatalLevel logs message and os.Exit(1)
	FatalLevel Level = zap.FatalLevel // 5

	DebugLevel Level = zap.DebugLevel // -1
)

type Field = zap.Field

type Logger struct {
	l     *zap.Logger // thread-safe
	level Level
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.l.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.l.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.l.Error(msg, fields...)
}
func (l *Logger) DPanic(msg string, fields ...Field) {
	l.l.DPanic(msg, fields...)
}
func (l *Logger) Panic(msg string, fields ...Field) {
	l.l.Panic(msg, fields...)
}
func (l *Logger) Fatal(msg string, fields ...Field) {
	l.l.Fatal(msg, fields...)
}

// function variables for all field types
// in github.com/uber-go/zap/field.go

// Expose as package level functions
var (
	Skip       = zap.Skip
	Binary     = zap.Binary
	Bool       = zap.Bool
	Boolp      = zap.Boolp
	ByteString = zap.ByteString

	Int    = zap.Int
	String = zap.String

	Float64   = zap.Float64
	Float64p  = zap.Float64p
	Float32   = zap.Float32
	Float32p  = zap.Float32p
	Durationp = zap.Durationp

	Any = zap.Any

	// Exposed logging functions
	Info   = std.Info
	Warn   = std.Warn
	Error  = std.Error
	DPanic = std.DPanic
	Panic  = std.Panic
	Fatal  = std.Fatal
	Debug  = std.Debug
)

// Expose custom encoder function
// var (
// 	WithCaller    = zap.WithCaller
// 	AddStacktrace = zap.AddStacktrace
// )

// not safe for concurrent use
func ResetDefault(l *Logger) {
	std = l
	Info = std.Info
	Warn = std.Warn
	Error = std.Error
	DPanic = std.DPanic
	Panic = std.Panic
	Fatal = std.Fatal
	Debug = std.Debug
}

var std = New(os.Stderr, InfoLevel)

func Default() *Logger {
	return std
}

// zap.Option.apply is not exposable
// type Option = zap.Option

type Option interface {
	Apply(*Logger)
}

type zapoption struct {
	opts []zap.Option
}

func (z *zapoption) Apply(l *Logger) {}
func (z *zapoption) zapOptions() []zap.Option {
	return z.opts
}

func fromZapOption(opts ...zap.Option) *zapoption {
	return &zapoption{opts}
}

func WithCaller(enable bool) Option {
	return fromZapOption(
		zap.WithCaller(enable),
		zap.AddCallerSkip(1),
	)
}
func AddStacktrace(skip int) Option {
	return fromZapOption(zap.AddCallerSkip(skip))
}

type EncoderOption interface {
	Config(*zap.Config)
}

// convert from Config to EncoderOption
type EncoderConfigurer func(*zap.Config)

type encoderoption struct {
	f EncoderConfigurer
}

func (opt *encoderoption) Config(config *zap.Config) {
	opt.f(config)
}

// only for passing as Option
func (opt *encoderoption) Apply(*Logger) {}

func WithEncoderConfig(fn EncoderConfigurer) Option {
	return &encoderoption{fn}
}

type TimeFormater func(time.Time) string

type configtimeformat struct {
	fmt TimeFormater
}

func (c *configtimeformat) Apply(l *Logger) {}
func (c *configtimeformat) Config(config *zap.Config) {
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(c.fmt(t))
	}
}

func WithTimeFormater(enc TimeFormater) Option {
	return &configtimeformat{enc}
}
func WithTimeFormat(format string) Option {
	return WithTimeFormater(func(t time.Time) string { return t.Format(format) })
}

// Helper
func configEncoder(config *zap.Config, opts ...Option) {
	for _, configurer := range opts {
		if c, ok := configurer.(EncoderOption); ok {
			c.Config(config)
		}
	}
}

// Helper
func zapOptions(opts []Option) []zap.Option {
	// convert to zap.Option
	zapOpts := []zap.Option{}
	for _, opt := range opts {
		if zopt, ok := opt.(*zapoption); ok {
			zapOpts = append(zapOpts, zopt.zapOptions()...)
		}
	}
	return zapOpts
}

// New create a new logger (not support log rotating).
func New(writer io.Writer, level Level, opts ...Option) *Logger {
	if writer == nil {
		panic("the writer is nil")
	}
	cfg := zap.NewProductionConfig()
	// config
	configEncoder(&cfg, opts...)

	// cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	// 	enc.AppendString()
	// }

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.AddSync(writer),
		zapcore.Level(level),
	)
	// convert to zap.Option
	zapOpts := zapOptions(opts)
	logger := &Logger{
		l:     zap.New(core, zapOpts...),
		level: level,
	}
	return logger
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}

func Sync() error {
	if std != nil {
		return std.Sync()
	}
	return nil
}
