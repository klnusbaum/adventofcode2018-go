package main

import (
	"fmt"
	"io"
	"os"

	"github.com/klnusbaum/adventofcode2018/ll"
)

type histogram map[rune]int

func main() {
	filename := os.Args[1]

	checksummer := newChecksummer(filename)

	if err := checksummer.Check(os.Stdout); err != nil {
		fmt.Errorf("Couldn't get checksum: %v", err)
		os.Exit(1)
	}
}

type checksummer struct {
	filename       string
	wordsWithTwo   int
	wordsWithThree int
}

func newChecksummer(filename string) *checksummer {
	return &checksummer{
		filename: filename,
	}
}

func (cs *checksummer) Check(out io.Writer) error {
	loader := ll.NewLineLoader(cs.filename)
	lines, err := loader.Load()
	if err != nil {
		return fmt.Errorf("Couldn't load lines: %v", err)
	}

	for _, line := range lines {
		cs.processLine(line)
	}

	checksum := cs.wordsWithTwo * cs.wordsWithThree
	fmt.Fprintf(out, "Checksum: %d\n", checksum)
	return nil
}

func (cs *checksummer) processLine(line string) {
	hist := calcHistogram(line)
	var hasTwos, hasThrees bool
	for _, count := range hist {
		if count == 2 {
			hasTwos = true
		} else if count == 3 {
			hasThrees = true
		}

		if hasTwos && hasThrees {
			break
		}
	}

	if hasTwos {
		cs.wordsWithTwo++
	}
	if hasThrees {
		cs.wordsWithThree++
	}
}

func calcHistogram(word string) histogram {
	hist := make(histogram)
	for _, char := range word {
		hist[char] += 1
	}
	return hist
}
