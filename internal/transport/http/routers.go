package http

import (
	"github.com/go-chi/chi/v5"

	"{{index .App "git"}}/internal/transport/http/middleware"
)

type RouterOptions struct {
	Logging bool
}

func RegisterHTTPRoutes(v1Router chi.Router, opts RouterOptions) chi.Router {
	r := chi.NewRouter()

	if opts.Logging {
		r.Use(middleware.Logging)
	}

	r.Get("/health", health)

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/", v1Router)
	})

	return r
}
