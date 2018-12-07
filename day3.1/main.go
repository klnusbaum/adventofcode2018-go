package main

/*
--- Day 3: No Matter How You Slice It ---
The Elves managed to locate the chimney-squeeze prototype fabric for Santa's suit (thanks to someone who helpfully wrote its box IDs on the wall of the warehouse in the middle of the night). Unfortunately, anomalies are still affecting them - nobody can even agree on how to cut the fabric.

The whole piece of fabric they're working on is a very large square - at least 1000 inches on each side.

Each Elf has made a claim about which area of fabric would be ideal for Santa's suit. All claims have an ID and consist of a single rectangle with edges parallel to the edges of the fabric. Each claim's rectangle is defined as follows:

The number of inches between the left edge of the fabric and the left edge of the rectangle.
The number of inches between the top edge of the fabric and the top edge of the rectangle.
The width of the rectangle in inches.
The height of the rectangle in inches.
A claim like #123 @ 3,2: 5x4 means that claim ID 123 specifies a rectangle 3 inches from the left edge, 2 inches from the top edge, 5 inches wide, and 4 inches tall. Visually, it claims the square inches of fabric represented by # (and ignores the square inches of fabric represented by .) in the diagram below:

...........
...........
...#####...
...#####...
...#####...
...#####...
...........
...........
...........
The problem is that many of the claims overlap, causing two or more claims to cover part of the same areas. For example, consider the following claims:

#1 @ 1,3: 4x4
#2 @ 3,1: 4x4
#3 @ 5,5: 2x2
Visually, these claim the following areas:

........
...2222.
...2222.
.11XX22.
.11XX22.
.111133.
.111133.
........
The four square inches marked with X are claimed by both 1 and 2. (Claim 3, while adjacent to the others, does not overlap either of them.)

If the Elves all proceed with their own plans, none of them will have enough fabric. How many square inches of fabric are within two or more claims?*/

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
