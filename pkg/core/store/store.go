package store

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DSN                   string // 数据库连接字符串
	MaxIdleConnections    int    // 最大空闲连接数
	MaxOpenConnections    int    // 最大打开连接数
	MaxConnectionLifeTime int    // 连接最大生命周期，单位秒
	MaxConnectionIdleTime int    // 连接最大空闲时间，单位秒
	LogLevel              int    // 日志配置,1:Silent,2:Error,3:Warn,4:Info
}

type gormLogWriter struct {
	logger zerolog.Logger
	level  logger.LogLevel
}

func (w gormLogWriter) Printf(format string, v ...interface{}) {
	switch w.level {
	case logger.Error:
		w.logger.Error().Msgf(format, v...)
	case logger.Warn:
		w.logger.Warn().Msgf(format, v...)
	case logger.Info:
		w.logger.Info().Msgf(format, v...)
	default:
		w.logger.Error().Msgf(format, v...)
	}
}

func MustNew(c *Config) *gorm.DB {
	logLevel := logger.Silent
	if c.LogLevel != 0 {
		logLevel = logger.LogLevel(c.LogLevel)
	}
	logWriter := &gormLogWriter{
		logger: log.Logger,
		level:  logLevel,
	}
	db, err := gorm.Open(postgres.Open(c.DSN), &gorm.Config{
		Logger: logger.New(logWriter, logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: false,
			Colorful:                  false,
		})})
	if err != nil {
		log.Fatal().Err(err).Msgf("open %s db failed", c.DSN)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msgf("get %s db failed", c.DSN)
	}
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(c.MaxIdleConnections)
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(c.MaxOpenConnections)
	// 设置连接最大生命周期
	sqlDB.SetConnMaxLifetime(time.Duration(c.MaxConnectionLifeTime) * time.Second)
	// 设置连接最大空闲时间
	sqlDB.SetConnMaxIdleTime(time.Duration(c.MaxConnectionIdleTime) * time.Second)
	return db
}
