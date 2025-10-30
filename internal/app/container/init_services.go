package container

import (
	dummyService "{{index .App "git"}}/internal/modules/dummy/service"
)

func (c *Container) initServices() {
	c.Services.DummyService = dummyService.New()
}
