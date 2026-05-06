package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ningzining/cove/pkg/core/middleware"
	"github.com/rs/zerolog/log"
)

type Server struct {
	engine *gin.Engine

	config *Config
	srv    *http.Server
}

func MustNewServer(cfg *Config) *Server {
	return NewServer(cfg)
}

func NewServer(cfg *Config) *Server {
	gin.SetMode(cfg.Mode)
	// 创建gin引擎
	engine := gin.New()
	// 安装中间件
	engine.Use(middleware.Recovery()).
		Use(middleware.RequestId()).
		Use(middleware.Logger()).
		Use(middleware.NoCache).
		Use(middleware.Cors).
		Use(middleware.Secure).
		Use(middleware.I18n())

	return &Server{config: cfg, engine: engine}
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}

func (s *Server) Start() error {
	if s.config.Mode == ProdMode {
		gin.SetMode(gin.ReleaseMode)
	}
	// 启动http server
	s.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		Handler: s.engine,
	}

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msgf("ListenAndServe fail")
		}
	}()
	log.Info().Msgf("server start at %s:%d", s.config.Host, s.config.Port)

	return nil
}

func (s *Server) Shutdown() error {
	// 等待中断信号优雅地关闭服务器（3 秒超时)。
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info().Msgf("Shutting down server ...")

	// 创建 ctx 用于通知服务器 goroutine, 它有 3 秒时间完成当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 通知服务器关闭
	if s.srv != nil {
		err := s.srv.Shutdown(ctx)
		if err != nil {
			log.Error().Err(err).Msgf("Shutdown server fail")
			return err
		}
	}

	return nil
}
