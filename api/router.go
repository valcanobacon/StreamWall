package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/valcanobacon/StreamWall/money"
	"github.com/valcanobacon/StreamWall/music"
)

func NewRouter(bank *money.Bank, durations music.FileDurations) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/sessions", func(r chi.Router) {
		r.Post("/", newSessionCreateHandler(bank))
		r.Route("/{sessionID}", func(r chi.Router) {
			r.Get("/", newSessionGetHandler(bank))
			r.Route("/songs", func(r chi.Router) {
				r.Get("/*", newSongHandler(bank, durations))
			})
		})
	})
	return r
}
