package v1

import "github.com/go-chi/chi/v5"

func Router(dummyRouter chi.Router) chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Mount("/dummy", dummyRouter)
	})

	return r
}
