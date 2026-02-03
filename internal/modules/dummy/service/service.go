package service

import (
	"context"

	"{{index .App "git"}}/internal/modules/dummy"
)

// Service implements dummy.Service interface.
type Service struct { {{- if index .Modules "postgres"}}
	repo dummy.Repository{{end}}
}

// New creates Service instance.
func New({{if index .Modules "postgres"}}
	repo dummy.Repository,{{end}}
) *Service {
	return &Service{ {{- if index .Modules "postgres"}}
		repo: repo,{{end}}
	}
}

// Create dummy and notify about its creation.
func (s *Service) Create(ctx context.Context, name string) (*dummy.Dummy, error) { {{- if index .Modules "postgres"}}
	return s.repo.Create(ctx, name){{else}}
	var res *dummy.Dummy

	res = &dummy.Dummy{
		ID:   100,
		Name: name,
	}

	return res, nil{{end}}
}

// List dummy objects.
func (s *Service) List(ctx context.Context, request dummy.ListRequest) ([]dummy.Dummy, error) { {{- if index .Modules "postgres"}}
	return s.repo.List(ctx, request){{else}}
	return []dummy.Dummy{
		{
			ID:   1,
			Name: "First Object",
		},
		{
			ID:   2,
			Name: "Second Object",
		},
	}, nil{{end}}
}
