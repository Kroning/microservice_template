package container

import (
	"context"
	"fmt"
	"net/http"

	"github.com/oklog/run"
	"go.uber.org/zap"

	"{{index .App "git"}}/pkg/logger"
)

func (c *Container) RunPublicHTTP(g *run.Group) {
	cfg := c.App.Cfg.HTTPServer
	publicServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           c.Routers.ChiHTTPRouters,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	g.Add(func() error {
		logger.Info(c.App.Ctx, "HTTP server is starting on ", zap.Uint("port", cfg.Port))
		return publicServer.ListenAndServe()
	}, func(err error) {
		logger.Info(c.App.Ctx, "HTTP server graceful shutdown started", zap.String("reason", err.Error()))

		ctx, cancel := context.WithTimeout(c.App.Ctx, cfg.GracefulShutdownTimeout)
		defer cancel()

		if err := publicServer.Shutdown(ctx); err != nil {
			logger.Error(c.App.Ctx, "HTTP server shutdown error", zap.Error(err))
		}
	})
}
