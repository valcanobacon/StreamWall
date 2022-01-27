package main

import (
	"fmt"
	"log"
	"net/http"

	//"github.com/go-chi/chi"
)

func main() {
	// configure the songs directory name and port
	const songsDir = "songs"
	const port = 8080

	// add a handler for the song files
	http.Handle("/", addHeaders(monitorTimeForCredits(http.FileServer(http.Dir(songsDir)))))
	fmt.Printf("Starting server on %v\n", port)
	log.Printf("Serving %s on HTTP port: %v\n", songsDir, port)

	// serve and log errors
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

// addHeaders will act as middleware to give us CORS support
func addHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}

func monitorTimeForCredits(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v", r.URL)
		h.ServeHTTP(w, r)
	}
}