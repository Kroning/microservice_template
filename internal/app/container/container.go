package container

import (
	"context"

	"{{index .App "git"}}/internal/app/config"
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
	// SomeService some.Service
}

type routers struct {
	/*HTTPRouters      []server.RouterHTTP
	GRPCServices []server.ServiceGRPC
	GraphqlResolver graphql.ResolverRoot*/
}

func (di *Container) RunApp() {
	di.initApp()
	di.stopApp()
}
