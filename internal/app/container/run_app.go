package container

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"
	"go.uber.org/zap"

	"{{index .App "git"}}/pkg/logger"
)

func (c *Container) appRun() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	var g run.Group

	g.Add(func() error {
		sig := <-signals
		logger.Info(c.App.Ctx, "got signal", zap.String("signal", sig.String()))

		return fmt.Errorf("interrupted by signal: %s", sig.String())
	}, func(error) {
		signal.Stop(signals)
		close(signals)
	})
{{if index .Modules "http_chi"}}
	c.RunHTTP(&g)
{{end}}
	if err := g.Run(); err != nil {
		logger.Fatal(c.App.Ctx, "exited with error", zap.Error(err))
	}
}

func (c *Container) RunHTTP(g *run.Group) {
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
