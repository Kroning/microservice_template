package http

import (
	"github.com/go-chi/chi/v5"

	"{{index .App "git"}}/internal/modules/dummy"
)

// Router returns HTTP routes for module.
func Router(
	dummyService dummy.Service,
) chi.Router {
	d := dummyHandlers{service: dummyService}
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", d.listTest)
		r.Post("/", d.createTest)
	})

	return r
}

type dummyHandlers struct {
	service dummy.Service
}

// DummyResponse represents dummy data.
// swagger:model DummyResponse
type DummyResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ErrorResponse provides error message and code.
// swagger:model ErrorResponse
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
