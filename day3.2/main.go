package main

/*
Amidst the chaos, you notice that exactly one claim doesn't overlap by even a single square inch of fabric with any other claim. If you can somehow draw attention to it, maybe the Elves will be able to make Santa's suit after all!

For example, in the claims above, only claim 3 is intact after all claims are made.

What is the ID of the only claim that doesn't overlap?

*/

import (
	"fmt"
	"github.com/klnusbaum/adventofcode2018/ll"
	"io"
	"os"
	"strconv"
	"strings"
)

type square []int

type fabric [][]square

type claim struct {
	x  int
	y  int
	h  int
	w  int
	id int
}

func main() {
	filename := os.Args[1]

	filler := &filler{
		filename: filename,
	}

	if err := filler.fill(os.Stdout); err != nil {
		fmt.Errorf("Couldn't get checksum: %v", err)
		os.Exit(1)
	}
}

type filler struct {
	filename string
}

func (f *filler) fill(out io.Writer) error {
	claims, err := loadClaims(f.filename)
	if err != nil {
		return fmt.Errorf("Couldn't load claims: %v", err)
	}

	sheet, err := processClaims(claims)
	if err != nil {
		return fmt.Errorf("Error processing claims: %v", err)
	}

	singleId := findSingleId(sheet)
	fmt.Fprintf(out, "Id with no overlap is: %d\n", singleId)

	return nil
}

func loadClaims(filename string) ([]claim, error) {
	loader := ll.NewLineLoader(filename)
	lines, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("Couldn't load lines: %v", err)
	}

	claims := make([]claim, 0, len(lines))
	for i, line := range lines {
		claim, err := loadClaim(line)
		if err != nil {
			return nil, fmt.Errorf("Error loading claim at line %d: %v", i, err)
		}
		claims = append(claims, claim)
	}

	return claims, nil
}

func loadClaim(line string) (claim, error) {
	fields := strings.Fields(line)
	xandy := strings.Split(strings.TrimSuffix(fields[2], ":"), ",")
	x, _ := strconv.Atoi(xandy[0])
	y, _ := strconv.Atoi(xandy[1])

	wandh := strings.Split(fields[3], "x")
	w, _ := strconv.Atoi(wandh[0])
	h, _ := strconv.Atoi(wandh[1])

	id, _ := strconv.Atoi(strings.TrimPrefix(fields[0], "#"))

	return claim{
		x:  x,
		y:  y,
		w:  w,
		h:  h,
		id: id,
	}, nil
}

func processClaims(claims []claim) (fabric, error) {
	sheet := newSheet()
	for _, claim := range claims {
		processClaim(claim, sheet)
	}
	return sheet, nil
}

func newSheet() fabric {
	sheet := make(fabric, 1000)
	for i := range sheet {
		sheet[i] = make([]square, 1000)
		for j := range sheet[i] {
			sheet[i][j] = make(square, 0)
		}
	}
	return sheet
}

func processClaim(claim claim, sheet fabric) {
	for i := claim.y; i < claim.y+claim.h; i++ {
		for j := claim.x; j < claim.x+claim.w; j++ {
			sheet[i][j] = append(sheet[i][j], claim.id)
		}
	}
}

func findSingleId(sheet fabric) int {
	possibleIds := make(map[int]struct{})
	blackList := make(map[int]struct{})
	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			if len(sheet[i][j]) == 1 {
				_, onBlackList := blackList[sheet[i][j][0]]
				if !onBlackList {
					possibleIds[sheet[i][j][0]] = struct{}{}
				}
			} else {
				for _, id := range sheet[i][j] {
					delete(possibleIds, id)
					blackList[id] = struct{}{}
				}
			}
		}
	}

	for id, _ := range possibleIds {
		return id
	}
	return -1
}
