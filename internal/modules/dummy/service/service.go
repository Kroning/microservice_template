package service

import (
	"context"

	"{{index .App "git"}}/internal/modules/dummy"
)

// Service implements dummy.Service interface.
type Service struct {
	//repo   dummy.Repository
}

// New creates Service instance.
func New() *Service {
	return &Service{
		//repo:   storage,
	}
}

// Create dummy in database and notify about its creation.
func (s *Service) Create(ctx context.Context, name string) (*dummy.Dummy, error) {
	var res *dummy.Dummy

	res = &dummy.Dummy{
		ID:   100,
		Name: name,
	}

	return res, nil
}

// List dummy objects from storage.
func (s *Service) List(ctx context.Context, request dummy.ListRequest) ([]dummy.Dummy, error) {
	return []dummy.Dummy{
		{
			ID:   1,
			Name: "First Object",
		},
		{
			ID:   2,
			Name: "Second Object",
		},
	}, nil
}
