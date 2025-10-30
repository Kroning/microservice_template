package container

import (
	"fmt"
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
	c.RunPublicHTTP(&g)
{{end}}
	if err := g.Run(); err != nil {
		logger.Fatal(c.App.Ctx, "exited with error", zap.Error(err))
	}
}
