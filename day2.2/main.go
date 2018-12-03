package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/klnusbaum/adventofcode2018/ll"
)

type histogram map[rune]int

func main() {
	filename := os.Args[1]

	finder := newFinder(filename)

	if err := finder.find(os.Stdout); err != nil {
		fmt.Errorf("Couldn't get checksum: %v", err)
		os.Exit(1)
	}
}

type finder struct {
	filename string
}

func newFinder(filename string) *finder {
	return &finder{
		filename: filename,
	}
}

func (cs *finder) find(out io.Writer) error {
	loader := ll.NewLineLoader(cs.filename)
	lines, err := loader.Load()
	if err != nil {
		return fmt.Errorf("Couldn't loade file %q: %v", cs.filename, err)
	}

	id1, id2, err := similarLines(lines)
	if err != nil {
		return err
	}

	inter := intersect(id1, id2)

	fmt.Println("The two lines that differ by only 1 character are:")
	fmt.Printf("\t%s\n", id1)
	fmt.Printf("\t%s\n", id2)
	fmt.Println("The characters between them that are the same are:")
	fmt.Printf("\t%s\n", inter)

	return nil
}

func similarLines(lines []string) (string, string, error) {
	found1, found2 := "", ""
	for i, line1 := range lines {
		for j, line2 := range lines {
			if i == j {
				continue
			}
			if utf8.RuneCount([]byte(line1)) != utf8.RuneCount([]byte(line2)) {
				return "", "", fmt.Errorf("Lines of different length: %q and %q", line1, line2)
			}
			dist := distance(line1, line2)
			if dist == 1 {
				found1 = line1
				found2 = line2
				break
			}
		}
	}

	if found1 == "" || found2 == "" {
		return "", "", fmt.Errorf("Could not find lines with only 1 letter different")
	}

	return found1, found2, nil
}

func distance(line1, line2 string) int {
	iter1 := newRuneIter(line1)
	iter2 := newRuneIter(line2)

	dist := 0
	for iter1.hasNext() {
		next1 := iter1.next()
		next2 := iter2.next()
		if next1 != next2 {
			dist++
		}
	}
	return dist
}

func intersect(line1, line2 string) string {
	iter1 := newRuneIter(line1)
	iter2 := newRuneIter(line2)

	var sb strings.Builder
	for iter1.hasNext() {
		next1 := iter1.next()
		next2 := iter2.next()
		if next1 == next2 {
			sb.WriteRune(next1)
		}
	}

	return sb.String()
}

func newRuneIter(str string) *runeIter {
	return &runeIter{
		str: str,
		i:   0,
	}
}

type runeIter struct {
	str string
	i   int
}

func (ri *runeIter) hasNext() bool {
	return ri.i < len(ri.str)
}

func (ri *runeIter) next() rune {
	current, width := utf8.DecodeRuneInString(ri.str[ri.i:])
	ri.i += width
	return current
}
