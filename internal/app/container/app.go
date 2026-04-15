package container

import (
	"context"

	"go.uber.org/zap"

	"{{index .App "git"}}/internal/app/config"{{if index .Modules "postgres"}}
	"{{index .App "git"}}/pkg/storage"{{end}}
)

type App struct {
	Ctx       context.Context
	Cfg       *config.Config{{if index .Modules "postgres"}}
	DB        storage.AbstractDB
	Repos     *Repositories{{end}}
	Services  *Services
	Logger    *zap.Logger{{if index .Modules "transport"}}
	Transport *Transport{{end}}
}

func BuildApp(
	ctx context.Context,
	cfg *config.Config,{{if index .Modules "postgres"}}
	db storage.AbstractDB,{{end}}
) *App { {{- if index .Modules "postgres"}}
	repos := NewRepositoryModule(db){{end}}
	services := NewServicesModule({{if index .Modules "postgres"}}
		repos,{{end}}
	){{if index .Modules "transport"}}
	transport := NewTransportModule(services, cfg){{end}}
	return &App{
		Ctx:       ctx,
		Cfg:       cfg,
{{- if index .Modules "postgres"}}
		DB:        db,
		Repos:     repos,
{{- end}}
		Services:  services,{{if index .Modules "transport"}}
		Transport: transport,{{end}}
	}
}
