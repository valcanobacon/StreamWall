package music

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func LoadDurations(fileGlob string) map[string]float64 {
	durations := map[string]float64{}

	files, err := filepath.Glob(fileGlob)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range files {
		dir, _ := path.Split(strings.Replace(p, "songs/", "/", 1))

		file, err := os.Open(p)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			row := scanner.Text()
			columns := strings.Split(row, " ")
			if duration, err := strconv.ParseFloat(columns[1], 64); err == nil {
				durations[dir+columns[0]] = duration
			} else {
				fmt.Println(err)
			}

		}

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}

	}

	return durations
}
