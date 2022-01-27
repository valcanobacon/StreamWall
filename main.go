package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"os"
	"bufio"
	"strconv"

	//"github.com/go-chi/chi"
)

var durations map[string]float64

func main() {
	// configure the songs directory name and port
	const songsDir = "songs"
	const port = 8080


	durationFiles, err := filepath.Glob("songs/*/durations.txt")
	if err != nil {
		log.Fatal(err)
	}

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
				durations[songDir + columns[0]] = duration
			} else {
				fmt.Println(err)
			}

		}

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}

		fmt.Println(durations)

	}

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
		duration := durations[r.URL.String()]
		satsPerSecond := 0.5
		cost := duration * satsPerSecond
		log.Printf("%v %g seconds at %g costs %g", r.URL, duration, satsPerSecond, cost)
		h.ServeHTTP(w, r)
	}
}