package handler

import (
	"link-storage-service/internal/middleware"
	"net/http"
	"github.com/go-chi/chi/v5"
)

func NewRouter(h *LinkHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)

	r.Route("/links", func(r chi.Router) {
		r.Post("/", h.CreateLink)
		r.Get("/", h.GetAllLinks)

		r.Route("/{code}", func(r chi.Router) {
			r.Get("/", h.GetByShortCode)
			r.Delete("/", h.DeleteLinks)
			r.Get("/stats", h.GetStats)
		})
	})

	return r
}
