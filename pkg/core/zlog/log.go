package zlog

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Dir        string // 日志存储目录
	Filename   string // 日志文件名前缀
	MaxSize    int    // 单个日志文件最大大小(MB)
	MaxBackups int    // 保留的最大日志文件数
	MaxAge     int    // 日志文件保留天数
	Compress   bool   // 是否压缩归档日志
	Level      string // 日志级别: debug/info/warn/error/fatal/panic
}

func Setup(config *Config) error {
	// 初始化日志目录
	if config.Dir != "" {
		err := os.MkdirAll(config.Dir, 0755)
		if err != nil {
			return errors.Wrap(err, "create log dir")
		}
		err = createFile(filepath.Join(config.Dir, config.Filename+".log"))
		if err != nil {
			return err
		}
	}
	// 初始化日志记录器
	if err := initZeroLogger(config); err != nil {
		return err
	}
	return nil
}

func initZeroLogger(config *Config) error {
	level, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		return errors.Wrap(err, "parse log level")
	}
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	zerolog.SetGlobalLevel(level)
	logger := zerolog.New(newWriter(config)).
		With().
		Timestamp().
		CallerWithSkipFrameCount(2).
		Logger()
	log.Logger = logger
	return nil
}

func createFile(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return errors.Wrap(err, "open log file")
	}
	defer file.Close()
	return nil
}
