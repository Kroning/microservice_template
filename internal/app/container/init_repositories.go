package container

import (
	dummyRepository "github.com/Kroning/example_service/internal/modules/dummy/repository"
)

func (c *Container) initRepositories() {
	c.Repositories.DummyRepository = dummyRepository.NewPostgresRepository(c.DB)
}
