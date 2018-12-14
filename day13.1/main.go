package main

// start Dec 14 5:30 am
// finish part 1: 8:08 am
/*
 */

import (
	"fmt"
	"github.com/klnusbaum/adventofcode2018/ll"
	"os"
	"sort"
)

type cart struct {
	direction rune
	lastTurn  string
}

type point struct {
	x int
	y int
}

type points []point

func (p points) Len() int      { return len(p) }
func (p points) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p points) Less(i, j int) bool {
	if p[i].y < p[j].y {
		return true
	} else if p[i].y == p[j].y && p[i].x < p[j].x {
		return true
	}

	return false
}

func main() {
	if err := analyze("input.txt"); err != nil {
		fmt.Printf("Error analyzing: %v\n", err)
		os.Exit(1)
	}
}

func analyze(filename string) error {
	loader := ll.NewLineLoader(filename)
	lines, err := loader.Load()
	if err != nil {
		return fmt.Errorf("Couldn't load lines: %v\n", err)
	}

	tracks, carts := loadTracksAndCarts(lines)
	// maxX, maxY := maxInt(tracks)
	// printBoard(tracks, carts, maxX, maxY)

	collisionPoint := runSimulation(tracks, carts)

	fmt.Printf("Collision at %d,%d\n", collisionPoint.x, collisionPoint.y)
	return nil
}

func runSimulation(tracks map[point]rune, carts map[point]cart) point {
	// maxX, maxY := maxInt(tracks)
	collision := false
	collisionPoint := point{}
	// printBoard(tracks, carts, maxX, maxY)
	for !collision {
		// for g := 0; g < 100 && !collision; g++ {
		ps := sortedCartPoints(carts)
		newCarts := make(map[point]cart)
		for _, p := range ps {
			currentCart, ok := carts[p]
			if !ok {
				panic(fmt.Sprintf("No cart at %d,%d!", p.x, p.y))
			}

			newPos := getNewPos(currentCart, p)

			if _, ok := carts[newPos]; ok {
				collision = true
				collisionPoint = newPos
				break
			}

			delete(carts, newPos)

			if _, ok := newCarts[newPos]; ok {
				collision = true
				collisionPoint = newPos
				break
			}

			trackAtNewPos := tracks[newPos]
			newDirection := currentCart.direction
			newLastTurn := currentCart.lastTurn
			if trackAtNewPos == '\\' {
				if currentCart.direction == '>' {
					newDirection = 'v'
				} else if currentCart.direction == '^' {
					newDirection = '<'
				} else if currentCart.direction == '<' {
					newDirection = '^'
				} else if currentCart.direction == 'v' {
					newDirection = '>'
				}
			} else if trackAtNewPos == '/' {
				if currentCart.direction == '>' {
					newDirection = '^'
				} else if currentCart.direction == '^' {
					newDirection = '>'
				} else if currentCart.direction == '<' {
					newDirection = 'v'
				} else if currentCart.direction == 'v' {
					newDirection = '<'
				}
			} else if trackAtNewPos == '+' {
				if currentCart.lastTurn == "right" {
					newDirection = turnLeft(currentCart.direction)
					newLastTurn = "left"
				} else if currentCart.lastTurn == "left" {
					newLastTurn = "straight"
				} else if currentCart.lastTurn == "straight" {
					newDirection = turnRight(currentCart.direction)
					newLastTurn = "right"
				}
			}

			newCarts[newPos] = cart{
				direction: newDirection,
				lastTurn:  newLastTurn,
			}
		}
		carts = newCarts
		// printBoard(tracks, carts, maxX, maxY)
	}

	// printBoard(tracks, carts, maxX, maxY)
	// for p, _ := range carts {
	// 	fmt.Printf("Cart at %d,%d\n", p.x, p.y)
	// }

	return collisionPoint
}

func getNewPos(currentCart cart, currentPoint point) point {
	if currentCart.direction == '<' {
		return point{currentPoint.x - 1, currentPoint.y}
	} else if currentCart.direction == '>' {
		return point{currentPoint.x + 1, currentPoint.y}
	} else if currentCart.direction == 'v' {
		return point{currentPoint.x, currentPoint.y + 1}
	} else {
		return point{currentPoint.x, currentPoint.y - 1}
	}
}

func turnLeft(r rune) rune {
	if r == '>' {
		return '^'
	} else if r == '<' {
		return 'v'
	} else if r == '^' {
		return '<'
	} else {
		return '>'
	}
}

func turnRight(r rune) rune {
	if r == '>' {
		return 'v'
	} else if r == '<' {
		return '^'
	} else if r == '^' {
		return '>'
	} else {
		return '<'
	}

}

func sortedCartPoints(carts map[point]cart) points {
	ps := make(points, 0, len(carts))
	for p, _ := range carts {
		ps = append(ps, p)
	}
	sort.Sort(ps)
	// fmt.Printf("%v\n", ps)

	return ps
}

func loadTracksAndCarts(lines []string) (map[point]rune, map[point]cart) {
	tracks := make(map[point]rune)
	carts := make(map[point]cart)
	for i, line := range lines {
		for j, r := range line {
			if r == ' ' {
				continue
			}
			p := point{j, i}

			if r == '-' || r == '|' || r == '/' || r == '\\' || r == '+' {
				tracks[p] = r
				continue
			}

			if r == '<' {
				tracks[p] = '-'
				carts[p] = cart{
					direction: r,
					lastTurn:  "right",
				}
			} else if r == '>' {
				tracks[p] = '-'
				carts[p] = cart{
					direction: r,
					lastTurn:  "right",
				}
			} else if r == 'v' {
				tracks[p] = '|'
				carts[p] = cart{
					direction: r,
					lastTurn:  "right",
				}
			} else if r == '^' {
				tracks[p] = '|'
				carts[p] = cart{
					direction: r,
					lastTurn:  "right",
				}
			}
		}
	}
	return tracks, carts
}

func printBoard(tracks map[point]rune, carts map[point]cart, maxX, maxY int) {
	for j := 0; j <= maxY; j++ {
		for i := 0; i <= maxX; i++ {
			if cart, ok := carts[point{i, j}]; ok {
				fmt.Printf("%c", cart.direction)
				continue
			}

			if track, ok := tracks[point{i, j}]; ok {
				fmt.Printf("%c", track)
				continue
			}

			fmt.Printf(" ")
		}
		fmt.Println()
	}
}

func maxInt(tracks map[point]rune) (int, int) {
	maxX := 0
	maxY := 0
	for p, _ := range tracks {
		if p.x > maxX {
			maxX = p.x
		}
		if p.y > maxY {
			maxY = p.y
		}
	}

	return maxX, maxY
}
