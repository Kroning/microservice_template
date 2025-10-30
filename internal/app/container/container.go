package container

import (
	"context"

	"github.com/go-chi/chi/v5"

	"{{index .App "git"}}/internal/app/config"
	"{{index .App "git"}}/internal/modules/dummy"
)

type Container struct {
	App             app
	Repositories    repositories
	InternalClients internalClients
	ExternalClients externalClients
	Services        services
	Routers         routers
}

type app struct {
	Ctx context.Context
	Cfg *config.Config
}

type repositories struct {
}

// internalClients contains internal clients (this company)
type internalClients struct {
	// internalClients someapi.Client
}

type externalClients struct {
	// GoogleAuthClient googleAuth.Client
}

type services struct {
	DummyService dummy.Service
}

type routers struct {
	ChiHTTPRouters chi.Router
}

func (c *Container) RunApp() {
	c.initApp()
	c.initServices()
{{- if index .Modules "http_chi"}}
	c.initRouters(){{end}}
	c.appRun()
	c.stopApp()
}
