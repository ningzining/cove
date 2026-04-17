package zlog

import (
	"io"
	"path/filepath"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func newWriter(cfg *Config) io.Writer {
	if cfg.Dir == "" {
		return zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = "2006-01-02 15:04:05.000"
		})
	}
	return &lumberjack.Logger{
		Filename:   filepath.Join(cfg.Dir, cfg.Filename+".log"),
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
}
