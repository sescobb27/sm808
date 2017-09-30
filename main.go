package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func trapExit(done chan<- struct{}) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT,
	)
	go func() {
		<-sigs
		fmt.Println("Terminating")
		done <- struct{}{}
		os.Exit(0)
	}()
}

func main() {
	done := make(chan struct{}, 1)
	trapExit(done)

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	defer file.Close()
	song := parseFile(file)

	ch := make(chan string)
	go song.Play(ch, done)
	printer(ch)
}

func parseFile(file io.Reader) *Song {
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	var text string

	song := &Song{Patterns: make([]*Pattern, 0)}
	for scanner.Scan() {
		text = scanner.Text()
		switch lineNumber {
		case 0:
			song.Name = text
		case 1:
			value, err := strconv.Atoi(text)
			if err != nil {
				log.Fatalf("invalid bpm format expected int: %v", err)
			}
			song.bpm = value
		default:
			inputPattern := strings.Split(text, " ")
			pattern := &Pattern{Name: inputPattern[0], Pattern: inputPattern[1]}
			song.Patterns = append(song.Patterns, pattern)
			if len(pattern.Pattern) > song.maxSteps {
				song.maxSteps = len(pattern.Pattern)
			}
		}
		lineNumber++
	}
	return song
}

func printer(read <-chan string) {
	for {
		select {
		case msg := <-read:
			fmt.Print(msg)
		}
	}
}

type Song struct {
	Name     string
	bpm      int
	Patterns []*Pattern
	maxSteps int
}
type Pattern struct {
	Name    string
	Pattern string
}

func (song *Song) Play(write chan<- string, done <-chan struct{}) {
	iter := 0
	output := make([]string, 0, song.maxSteps)

	for {
		select {
		case <-done:
			return
		default:
			for _, pattern := range song.Patterns {
				if string(pattern.Pattern[iter%len(pattern.Pattern)]) == "1" {
					output = append(output, pattern.Name)
				}
			}
			if len(output) == 0 {
				write <- "_"
			} else {
				write <- strings.Join(output, "+")
			}
			output = output[:0]
			iter++
		}
	}
}
