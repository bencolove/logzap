package logzap

import (
	"io"

	"gopkg.in/natefinch/lumberjack.v2"
)

type RotateOption struct {
	Filepath string
	// megabytes M
	MaxSize int
	// days
	MaxAge    int
	MaxBackup int
	// disabled by default
	Compress bool
}

func (opt *RotateOption) Writer() io.Writer {
	return NewRotater(opt)
}
func (opt *RotateOption) Apply(*Logger) {}

func NewRotater(opt *RotateOption) io.Writer {
	return &lumberjack.Logger{
		Filename:   opt.Filepath,
		MaxSize:    opt.MaxSize,
		MaxBackups: opt.MaxBackup,
		MaxAge:     opt.MaxAge,
		Compress:   opt.Compress,
	}
}
