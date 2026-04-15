package container

import (
	"context"
	"runtime"

	"{{index .App "git"}}/internal/app/config"

	"{{index .App "git"}}/pkg/logger"
)

func (c *Container) initApp() {
	c.Ctx = context.Background()
	c.Cfg = config.GetConfig(c.Ctx)

	logger.SetLevel(c.Cfg.App.LogLevel)

	logger.Info(c.Ctx, "Init app: platform: "+runtime.GOOS+"/"+runtime.GOARCH)
}
