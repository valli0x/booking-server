package httpserver

import "github.com/go-chi/chi"

func (s *server) routers() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Post("/orders", s.createOrder())
		r.Get("/orders/{userID}", s.getorders())
	})
	return r
}
