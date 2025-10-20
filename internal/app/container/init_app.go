package container

import (
	"context"
	"runtime"

	"{{index .App "git"}}/internal/app/config"

	"{{index .App "git"}}/pkg/logger"
)

func (di *Container) initApp() {
	di.App.Ctx = context.Background()
	di.App.Cfg = config.GetConfig(di.App.Ctx)

	logger.SetLevel(di.App.Cfg.App.LogLevel)

	logger.Info(di.App.Ctx, "Init app: platform: "+runtime.GOOS+"/"+runtime.GOARCH)
}
