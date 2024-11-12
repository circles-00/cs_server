package router

import (
	"net/http"

	"api/internal/handler"

	"github.com/go-chi/chi/v5"
)

func csRouter(handler *handler.CsHandler) http.Handler {
	r := chi.NewRouter()

	r.Post("/register", handler.RegisterServerHandler)

	return r
}