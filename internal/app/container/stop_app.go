package container

import (
	"{{index .App "git"}}/pkg/logger"
)

func (c *Container) stopApp() {
	logger.Info(c.App.Ctx, "Stopping app")

	_ = logger.Logger().Sync()
}
