package mw

import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

func Log(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			logger.Info("request",
				slog.String("path", c.Request.URL.Path),
			)
		}()
	}
}
