package main

import (
	"fmt"
	"github.com/klnusbaum/adventofcode2018/ll"
	"io"
	"os"
	"strconv"
	"strings"
)

type fabric [][]int

func main() {
	filename := os.Args[1]

	filler := newFiller(filename)

	if err := filler.fill(os.Stdout); err != nil {
		fmt.Errorf("Couldn't get checksum: %v", err)
		os.Exit(1)
	}
}

func newFiller(filename string) *filler {
	sheet := make(fabric, 1000)
	for i := range sheet {
		sheet[i] = make([]int, 1000)
	}

	return &filler{
		filename: filename,
		sheet:    sheet,
	}
}

type filler struct {
	filename string
	sheet    fabric
	claims   []claim
}

type claim struct {
	x int
	y int
	h int
	w int
}

func (f *filler) fill(out io.Writer) error {
	if err := f.loadClaims(); err != nil {
		return fmt.Errorf("Couldn't load claims: %v", err)
	}

	if err := f.processClaims(); err != nil {
		return fmt.Errorf("Couldn't process claims: %v", err)
	}

	numOverlap := f.overlaps()

	fmt.Fprintf(out, "Number of overlapping squares: %d\n", numOverlap)

	return nil
}

func (f *filler) loadClaims() error {
	loader := ll.NewLineLoader(f.filename)
	lines, err := loader.Load()

	if err != nil {
		return fmt.Errorf("Couldn't load lines: %v", err)
	}

	for i, line := range lines {
		if err := f.loadClaim(line); err != nil {
			return fmt.Errorf("Couldn't load claim on line %d: %v", i, err)
		}
	}

	return nil
}

func (f *filler) loadClaim(line string) error {
	fields := strings.Fields(line)
	xandy := strings.Split(strings.TrimSuffix(fields[2], ":"), ",")
	x, _ := strconv.Atoi(xandy[0])
	y, _ := strconv.Atoi(xandy[1])

	wandh := strings.Split(fields[3], "x")
	w, _ := strconv.Atoi(wandh[0])
	h, _ := strconv.Atoi(wandh[1])

	f.claims = append(f.claims, claim{
		x: x,
		y: y,
		w: w,
		h: h,
	})

	return nil
}

func (f *filler) processClaims() error {
	for _, claim := range f.claims {
		f.processClaim(claim)
	}
	return nil
}

func (f *filler) processClaim(claim claim) error {
	for i := claim.y; i < claim.y+claim.h; i++ {
		for j := claim.x; j < claim.x+claim.w; j++ {
			f.sheet[i][j] += 1
		}
	}

	return nil
}

func (f *filler) overlaps() int {
	numOverlaps := 0
	for _, line := range f.sheet {
		for _, square := range line {
			if square > 1 {
				numOverlaps++
			}
		}
	}

	return numOverlaps
}
