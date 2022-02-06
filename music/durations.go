package music

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type FileDurations = map[string]float64

func LoadDurations(glob, prefix string) (map[string]float64, error) {
	// filename -> durationInSeconds
	durations := FileDurations{}

	files, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}

	for _, p := range files {
		dir, _ := path.Split(strings.Replace(p, prefix, "/", 1))

		file, err := os.Open(p)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			row := scanner.Text()
			columns := strings.Split(row, " ")
			if d, err := strconv.ParseFloat(columns[1], 64); err == nil {
				durations[dir+columns[0]] = d
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}

	}

	return durations, nil
}
