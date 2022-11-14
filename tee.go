package logzap

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Tee
type LevelEnabler func(Level) bool

// One tee one core, one core one target one filter
type TeeOption struct {
	Writer       io.Writer
	LevelEnabler LevelEnabler
}

func NewTee(tees []TeeOption, opts ...Option) *Logger {
	var cores []zapcore.Core

	cfg := zap.NewProductionConfig()
	// config
	configEncoder(&cfg, opts...)

	for _, tee := range tees {
		if tee.Writer == nil {
			Panic("Can not use nil writer")
		}

		lf := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return tee.LevelEnabler(Level(lvl))
		})

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.AddSync(tee.Writer),
			lf,
		)

		cores = append(cores, core)
	}

	// convert to zap.Option
	zapOpts := zapOptions(opts)

	return &Logger{
		l: zap.New(zapcore.NewTee(cores...), zapOpts...),
	}
}
