package router

import (
	"api/internal/handler"
	"api/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func NewRouter(csHandler *handler.CsHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Mount("/", csRouter(csHandler))

	return r
}
