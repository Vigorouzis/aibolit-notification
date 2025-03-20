package http

import (
	"context"

	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/vigorouzis/aibolit-notification/internal/core/application"
	"github.com/vigorouzis/aibolit-notification/internal/interface/http/mw"
)

type Server struct {
	cfg    Config
	log    *slog.Logger
	router *gin.Engine
	srv    *application.ScheduleService
}

func New(cfg Config, app *application.ScheduleService, log *slog.Logger) (*Server, error) {
	s := &Server{
		log:    log,
		cfg:    cfg,
		srv:    app,
		router: gin.New(),
	}

	s.router.Use(mw.Log(s.log))

	s.router.POST("/schedule", s.createSchedule)
	s.router.GET("/schedules", s.getSchedulesByUserID)
	s.router.GET("/schedule", s.getSchedule)
	s.router.GET("/next_takings", s.nextTakings)

	return s, nil
}

func (s *Server) Run(ctx context.Context) error {
	s.router.Use(mw.Log(s.log))

	select {
	case <-ctx.Done():
		return fmt.Errorf("context done")
	default:
		return s.router.Run(s.cfg.Addr())
	}
}
