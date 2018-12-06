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
	"os"
	// "sort"
	"strconv"
	"strings"
	// "github.com/klnusbaum/adventofcode2018/ll"
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
	// loader := ll.NewLineLoader(filename)
	// lines, err := loader.Load()
	// if err != nil {
	// 	return fmt.Errorf("Error loading lines: %v", lines)
	// }

	// allPoints := getPoints(lines)
	// boundPoints := pointsWithBounds(allPoints)
	allPoints := map[int]point{
		0: point{1, 1},
		1: point{1, 6},
		2: point{8, 3},
		3: point{3, 4},
		4: point{5, 5},
		5: point{8, 9},
	}
	boundPoints := pointsWithBounds(allPoints)

	// var bp1 []int
	// for id, _ := range boundPoints {
	// 	bp1 = append(bp1, id)
	// }
	// sort.Ints(bp1)
	// fmt.Printf("Bound points: %v\n", bp1)

	closestHisto := make(map[int]int)
	for i := -100; i < 900; i++ {
		for j := -100; j < 900; j++ {
			closestId := closestPoint(allPoints, i, j)
			if closestId == -1 {
				continue
			}
			if _, ok := boundPoints[closestId]; ok {
				closestHisto[closestId] += 1
			}
		}
	}

	fmt.Printf("Closest histo: %v\n", closestHisto)

	largestArea := largestInHisto(closestHisto)

	fmt.Printf("Largets bounded area: %v\n", largestArea)
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
		if isBounded(points, point) {
			boundedPoints[id] = point
		}
	}
	return boundedPoints
}

func isBounded(allPoints map[int]point, target point) bool {
	above, below, left, right := false, false, false, false
	for _, p := range allPoints {
		if p.x < target.x {
			left = true
		} else if p.x > target.x {
			right = true
		}

		if p.y < target.y {
			above = true
		} else if p.y > target.y {
			below = true
		}

		if above && below && left && right {
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

func largestInHisto(histo map[int]int) int {
	largestArea := 0
	for _, area := range histo {
		if area > largestArea {
			largestArea = area
		}
	}
	return largestArea

}
