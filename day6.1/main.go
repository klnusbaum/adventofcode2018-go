package main

/*
The device on your wrist beeps several times, and once again you feel like you're falling.

"Situation critical," the device announces. "Destination indeterminate. Chronal interference detected. Please specify new target coordinates."

The device then produces a list of coordinates (your puzzle input). Are they places it thinks are safe or dangerous? It recommends you check manual page 729. The Elves did not give you a manual.

If they're dangerous, maybe you can minimize the danger by finding the coordinate that gives the largest distance from the other points.

Using only the Manhattan distance, determine the area around each coordinate by counting the number of integer X,Y locations that are closest to that coordinate (and aren't tied in distance to any other coordinate).

Your goal is to find the size of the largest area that isn't infinite. For example, consider the following list of coordinates:

1, 1
1, 6
8, 3
3, 4
5, 5
8, 9
If we name these coordinates A through F, we can draw them on a grid, putting 0,0 at the top left:

..........
.A........
..........
........C.
...D......
.....E....
.B........
..........
..........
........F.
This view is partial - the actual grid extends infinitely in all directions. Using the Manhattan distance, each location's closest coordinate can be determined, shown here in lowercase:

aaaaa.cccc
aAaaa.cccc
aaaddecccc
aadddeccCc
..dDdeeccc
bb.deEeecc
bBb.eeee..
bbb.eeefff
bbb.eeffff
bbb.ffffFf
Locations shown as . are equally far from two or more coordinates, and so they don't count as being closest to any.

In this example, the areas of coordinates A, B, C, and F are infinite - while not shown here, their areas extend forever outside the visible grid. However, the areas of coordinates D and E are finite: D is closest to 9 locations, and E is closest to 17 (both including the coordinate's location itself). Therefore, in this example, the size of the largest area is 17.

What is the size of the largest area that isn't infinite?
*/

import (
	"fmt"
	"github.com/klnusbaum/adventofcode2018/ll"
	"os"
	"sort"
	"strconv"
	"strings"
)

type point struct {
	x int
	y int
}

func main() {
	if err := analyze("input.txt"); err != nil {
		fmt.Printf("Error analyzing schedules: %v\n", err)
		os.Exit(1)
	}
}

func analyze(filename string) error {
	loader := ll.NewLineLoader(filename)
	lines, err := loader.Load()
	if err != nil {
		return fmt.Errorf("Error loading lines: %v", lines)
	}

	allPoints := getPoints(lines)
	boundPoints := pointsWithBounds(allPoints)

	// allPoints := map[int]point{
	// 	0: point{1, 1},
	// 	1: point{1, 6},
	// 	2: point{8, 3},
	// 	3: point{3, 4},
	// 	4: point{5, 5},
	// 	5: point{8, 9},
	// }
	// boundPoints := pointsWithBounds(allPoints)

	var bp1 []int
	for id, _ := range boundPoints {
		bp1 = append(bp1, id)
	}
	sort.Ints(bp1)
	fmt.Printf("Bound points: %v\n", bp1)

	// maxX, maxY := maxXY(allPoints)

	// fmt.Printf("Maxx %d Maxy %d\n", maxX, maxY)

	closestHisto := make(map[int]int)
	for i := 0; i < 500; i++ {
		for j := 0; j < 500; j++ {
			closestId := closestPoint(allPoints, i, j)
			if closestId == -1 {
				continue
			}
			if _, ok := boundPoints[closestId]; ok {
				closestHisto[closestId] += 1
			}
		}
	}

	fmt.Printf("bound points: %v\n", boundPoints)
	fmt.Printf("Closest histo: %v\n", closestHisto)

	id, largestArea := largestInHisto(closestHisto)

	fmt.Printf("Largest bounded area is point %d with area: %d\n", id, largestArea)
	return nil
}

func getPoints(lines []string) map[int]point {
	points := make(map[int]point)
	for i, line := range lines {
		fields := strings.Split(line, ", ")
		x, _ := strconv.Atoi(fields[0])
		y, _ := strconv.Atoi(fields[1])
		points[i] = point{x, y}
	}

	return points
}

func pointsWithBounds(points map[int]point) map[int]point {
	boundedPoints := make(map[int]point)
	for id, point := range points {
		fmt.Printf("Checking Point %d: %d, %d\n", id, point.x, point.y)
		if isBounded(points, point) {
			boundedPoints[id] = point
		}
	}
	return boundedPoints
}

func isBounded(allPoints map[int]point, target point) bool {
	topRight, topLeft, bottomRight, bottomLeft := false, false, false, false

	for _, p := range allPoints {
		if p.x == target.x && p.y == target.y {
			continue
		}

		if p.x > target.x && p.y > target.y {
			fmt.Printf("\tTop right bounded by %d, %d\n", p.x, p.y)
			topRight = true
		} else if p.x < target.x && p.y > target.y {
			fmt.Printf("\tTop left bounded by %d, %d\n", p.x, p.y)
			topLeft = true
		} else if p.x > target.x && p.y < target.y {
			fmt.Printf("\tBottom right bounded by %d, %d\n", p.x, p.y)
			bottomRight = true
		} else if p.x < target.x && p.y < target.y {
			fmt.Printf("\tBottom left bounded by %d, %d\n", p.x, p.y)
			bottomLeft = true
		}

		if topRight && topLeft && bottomRight && bottomLeft {
			return true
		}

	}

	return false

}

func closestPoint(allPoints map[int]point, x, y int) int {
	dists := make(map[int][]int)
	for id, point := range allPoints {
		dist := manDistance(point.x, point.y, x, y)
		dists[dist] = append(dists[dist], id)
	}

	shortestDistance := 100000000
	var shortestIds []int
	for dist, ids := range dists {
		if dist < shortestDistance {
			shortestDistance = dist
			shortestIds = ids
		}
	}

	if len(shortestIds) == 1 {
		return shortestIds[0]
	}

	return -1
}

func manDistance(x1, y1, x2, y2 int) int {
	return abs(x1-x2) + abs(y1-y2)
}

func abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
}

func largestInHisto(histo map[int]int) (int, int) {
	largestArea := 0
	largestId := -1
	for id, area := range histo {
		if area > largestArea {
			largestArea = area
			largestId = id
		}
	}
	return largestId, largestArea

}

func maxXY(points map[int]point) (int, int) {
	xs := make([]int, 0, len(points))
	ys := make([]int, 0, len(points))

	for _, point := range points {
		xs = append(xs, point.x)
		ys = append(ys, point.y)
	}

	maxX := xs[0]
	for i := 1; i < len(xs); i++ {
		if xs[i] > maxX {
			maxX = xs[i]
		}
	}

	maxY := ys[0]
	for i := 1; i < len(ys); i++ {
		if ys[i] > maxY {
			maxY = ys[i]
		}
	}

	return maxX, maxY
}

// bad is bounded
// func isBounded(allPoints map[int]point, target point) bool {
// 	above, below, left, right := false, false, false, false
// 	topPoint := point{
// 		x: target.x,
// 		y: 10000000,
// 	}
// 	bottomPoint := point{
// 		x: target.x,
// 		y: -10000000,
// 	}
// 	leftPoint := point{
// 		x: -10000000,
// 		y: target.y,
// 	}
// 	rightPoint := point{
// 		x: 10000000,
// 		y: target.y,
// 	}

// 	for _, p := range allPoints {
// 		if p.x == target.x && p.y == target.y {
// 			continue
// 		}

// 		tpDist := manDistance(topPoint.x, topPoint.y, p.x, p.y)
// 		bpDist := manDistance(bottomPoint.x, bottomPoint.y, p.x, p.y)
// 		lpDist := manDistance(leftPoint.x, leftPoint.y, p.x, p.y)
// 		rpDist := manDistance(rightPoint.x, rightPoint.y, p.x, p.y)
// 		if tpDist < manDistance(topPoint.x, topPoint.y, target.x, target.y) {
// 			above = true
// 		}

// 		if bpDist < manDistance(bottomPoint.x, bottomPoint.y, target.x, target.y) {
// 			below = true
// 		}

// 		if lpDist < manDistance(leftPoint.x, leftPoint.y, target.x, target.y) {
// 			left = true
// 		}

// 		if rpDist < manDistance(rightPoint.x, rightPoint.y, target.x, target.y) {
// 			left = true
// 		}

// 		if above && below && left && right {
// 			return true
// 		}
// 	}

// 	return false
// }
