package container

import (
	"context"

	"{{index .App "git"}}/internal/app/config"{{if index .Modules "postgres"}}
	"{{index .App "git"}}/pkg/storage"{{end}}
)

type Container struct {
	Ctx             context.Context
	Cfg             *config.Config
	App             *App{{if index .Modules "postgres"}}
	DB              storage.AbstractDB{{end}}
	InternalClients internalClients
	ExternalClients externalClients
}

// internalClients contains internal clients (this company)
type internalClients struct {
	// internalClients someapi.Client
}

type externalClients struct {
	// GoogleAuthClient googleAuth.Client
}

func (c *Container) RunApp() {
	c.initApp(){{if index .Modules "postgres"}}
	c.initDatabase(){{end}}

	c.App = BuildApp(
		c.Ctx,
		c.Cfg,{{if index .Modules "postgres"}}
		c.DB,{{end}}
	)

	c.appRun()
	c.stopApp()
}
