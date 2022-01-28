package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

var durations map[string]float64

var sessions map[string]*Session

func main() {
	// configure the songs directory name and port
	const songsDir = "songs"
	const port = 8080

	durationFiles, err := filepath.Glob("songs/*/durations.txt")
	if err != nil {
		log.Fatal(err)
	}

	sessions = map[string]*Session{}
	durations = map[string]float64{}

	for _, durationFilePath := range durationFiles {
		durationFilePathWithoutPrefix := strings.Replace(durationFilePath, "songs/", "/", 1)
		songDir, _ := path.Split(durationFilePathWithoutPrefix)
		fmt.Println(songDir)

		file, err := os.Open(durationFilePath)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			row := scanner.Text()
			columns := strings.Split(row, " ")
			if duration, err := strconv.ParseFloat(columns[1], 64); err == nil {
				durations[songDir+columns[0]] = duration
			} else {
				fmt.Println(err)
			}

		}

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}

		fmt.Println(durations)

	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/sessions", func(r chi.Router) {
		r.Post("/", sessionCreateHandler)
		r.Route("/{sessionID}", func(r chi.Router) {
			r.Get("/", sessionGetHandler)
			r.Route("/streams", func(r chi.Router) {
				r.Get("/*", streamHandler)
			})
		})
	})

	fmt.Printf("Starting server on %v\n", port)
	log.Printf("Serving %s on HTTP port: %v\n", songsDir, port)

	// serve and log errors
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	prefix := "/sessions/" + sessionID + "/streams"
	stream := strings.Replace(r.URL.String(), prefix, "", 1)

	if strings.HasSuffix(r.URL.String(), ".ts") {

		_, ok := sessions[sessionID]
		if !ok {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		duration := durations[stream]
		satsPerSecond := 0.5
		cost := duration * satsPerSecond
		log.Printf("%v %g seconds at %g costs %g", r.URL, duration, satsPerSecond, cost)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fileServer := http.FileServer(http.Dir("songs"))
	h := http.StripPrefix(prefix, fileServer)
	h.ServeHTTP(w, r)
}

func sessionCreateHandler(w http.ResponseWriter, r *http.Request) {
	newId := uuid.New()
	session := NewSession(&newId, 0)
	sessions[session.Id.String()] = session
	render.Render(w, r, session)
}

func sessionGetHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "sessionID")
	if session, ok := sessions[sessionID]; ok {
		render.Render(w, r, session)
		return
	}
	http.Error(w, http.StatusText(404), 404)
}

type Session struct {
	Id      *uuid.UUID `json:"id"`
	Credits int        `json:"credits"`
}

func NewSession(id *uuid.UUID, credits int) *Session {
	response := &Session{Id: id, Credits: credits}
	return response
}

func (s *Session) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
