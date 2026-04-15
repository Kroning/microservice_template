package container

import (
	"{{index .App "git"}}/internal/modules/dummy"
	dummyRepository "{{index .App "git"}}/internal/modules/dummy/repository"
	"{{index .App "git"}}/pkg/storage"
)

type Repositories struct {
	DummyRepository dummy.Repository
}

func NewRepositoryModule(db storage.AbstractDB) *Repositories {
	return &Repositories{
		dummyRepository.NewPostgresRepository(db),
	}
}
