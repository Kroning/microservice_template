package container

import (
	dummy "{{index .App "git"}}/internal/modules/dummy/transport/http"
	"{{index .App "git"}}/internal/transport/http"
	"{{index .App "git"}}/internal/transport/http/v1"
)

func (c *Container) initRouters() {
	c.initHTTPRouters()
}

func (c *Container) initHTTPRouters() {
	dummyRouter := dummy.Router(c.Services.DummyService)

	v1Router := v1.Router(dummyRouter)

	c.Routers.ChiHTTPRouters = http.RegisterHTTPRoutes(v1Router)
}
