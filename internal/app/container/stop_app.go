package container

import (
	"{{index .App "git"}}/pkg/logger"
)

func (di *Container) stopApp() {
	logger.Info(di.App.Ctx, "Stopping app")

	_ = logger.Logger().Sync()
}
