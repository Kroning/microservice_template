package container

import (
	"{{index .App "git"}}/internal/modules/dummy"
	dummyService "{{index .App "git"}}/internal/modules/dummy/service"
)

type Services struct {
	DummyService dummy.Service
}

func NewServicesModule(
{{- if index .Modules "postgres"}}
	repos *Repositories,
{{- end}}
	) *Services {
	return &Services{
		DummyService: dummyService.New(
{{- if index .Modules "postgres"}}
			repos.DummyRepository,
{{- end}}
		),
	}
}
