package http

import (
	"github.com/go-chi/chi/v5"
)

func RegisterHTTPRoutes(v1Router chi.Router) chi.Router {
	r := chi.NewRouter()

	r.Get("/health", health)

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/", v1Router)
	})

	return r
}
