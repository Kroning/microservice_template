package container

import ({{if index .Modules "postgres"}}
	"go.uber.org/zap"
{{end}}
	"{{index .App "git"}}/pkg/logger"
)

func (c *Container) stopApp() {
	logger.Info(c.App.Ctx, "Stopping app")
{{if index .Modules "postgres"}}
	if c.DB != nil {
		if err := c.DB.Close(); err != nil {
			logger.Error(c.App.Ctx, "failed to close database connection", zap.Error(err))
		} else {
			logger.Info(c.App.Ctx, "database connection closed")
		}
	}
{{end}}
	_ = logger.Logger().Sync()
}
