package rest

import (
	"fmt"
	"net/http"

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
	log.Info().Msgf("server start at %s:%d", s.config.Host, s.config.Port)

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msgf("ListenAndServe fail")
	}

	return nil
}
