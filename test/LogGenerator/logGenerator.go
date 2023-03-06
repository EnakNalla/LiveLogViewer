package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-loremipsum/loremipsum"
)

func main() {
	path := flag.String("path", "", "Path to log file")
	levels := flag.Bool("levels", false, "Include log levels")
	interval := flag.Int("interval", 1000, "Time between writes in ms")
	dateTime := flag.Bool("date", true, "Include dateTime")
	flag.Parse()

	if *path == "" {
		fmt.Println("Path is required")
		os.Exit(1)
	}

	fileHandle, err := os.OpenFile(*path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Writing to file: " + *path + " with interval: " + fmt.Sprint(*interval) + "ms")

	for {
		time.Sleep(time.Duration(*interval) * time.Millisecond)
		_, err = fileHandle.WriteString(buildLine(loremipsum.New().Sentence(), *dateTime, *levels))
	}
}

var Levels = []string{"[INF]", "[DBG]", "[ERR]", "[WARN]"}

func buildLine(sentence string, date bool, levels bool) string {
	line := ""

	if date {
		line += "[" + time.Now().Format("2006-01-02 15:04:05.000") + "]" + " "
	}

	if levels {
		line += Levels[rand.Intn(4)] + " "
	}

	line += sentence + "\n"

	return line
}
