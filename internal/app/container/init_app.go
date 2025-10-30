package container

import (
	"context"
	"runtime"

	"{{index .App "git"}}/internal/app/config"

	"{{index .App "git"}}/pkg/logger"
)

func (c *Container) initApp() {
	c.App.Ctx = context.Background()
	c.App.Cfg = config.GetConfig(c.App.Ctx)

	logger.SetLevel(c.App.Cfg.App.LogLevel)

	logger.Info(c.App.Ctx, "Init app: platform: "+runtime.GOOS+"/"+runtime.GOARCH)
}
