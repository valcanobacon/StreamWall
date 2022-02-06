package api

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/valcanobacon/StreamWall/money"
	"github.com/valcanobacon/StreamWall/music"
)

func newSongHandler(bank *money.Bank, durations music.FileDurations) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := chi.URLParam(r, "sessionID")
		prefix := "/sessions/" + sessionID + "/songs"
		stream := strings.Replace(r.URL.String(), prefix, "", 1)

		if strings.HasSuffix(r.URL.String(), ".ts") {

			sid, err := uuid.Parse(chi.URLParam(r, "sessionID"))
			if err != nil {
				return
			}

			s := bank.GetSession(sid)
			if s == nil {
				http.Error(w, http.StatusText(404), 404)
				return
			}

			duration := durations[stream]
			satsPerSecond := 1.0
			cost := int64(duration * satsPerSecond)

			if s.Credits < cost {
				http.Error(w, http.StatusText(403), 403)
				return
			}

			// log.Printf(
			// 	"%v %f seconds at %g costs %g", r.URL, duration, satsPerSecond, cost,
			// )

			bank.NewTransaction(sid, (-1 * cost))
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")

		fileServer := http.FileServer(http.Dir("songs"))
		h := http.StripPrefix(prefix, fileServer)
		h.ServeHTTP(w, r)
	}
}
