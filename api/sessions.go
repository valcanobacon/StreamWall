package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/valcanobacon/StreamWall/money"
)

func newSessionCreateHandler(bank *money.Bank) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s := bank.NewSession()
		render.Render(w, r, s)
	}
}

func newSessionGetHandler(bank *money.Bank) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sid, err := uuid.Parse(chi.URLParam(r, "sessionID"))
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
		}
		s := bank.GetSession(sid)
		if s == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		render.Render(w, r, s)
	}
}
