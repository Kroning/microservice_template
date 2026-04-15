package container

import (
	"github.com/go-chi/chi/v5"

	"{{index .App "git"}}/internal/app/config"
	dummy "{{index .App "git"}}/internal/modules/dummy/transport/http"
	"{{index .App "git"}}/internal/transport/http"
	"{{index .App "git"}}/internal/transport/http/v1"
)

type Transport struct {
	ChiRouter chi.Router
}

func NewTransportModule(services *Services, cfg *config.Config) *Transport {
	dummyRouter := dummy.Router(services.DummyService)
	v1Router := v1.Router(dummyRouter)
	return &Transport{
		ChiRouter: http.RegisterHTTPRoutes(v1Router, http.RouterOptions{
			Logging: cfg.HTTPServer.Logging,
		}),
	}
}
