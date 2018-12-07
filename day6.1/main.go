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
	maxX, maxY := maxXY(allPoints)

	closestHisto := make(map[int]int)
	pointOwnership := make(map[int]map[int]int)
	for i := 0; i < maxX; i++ {
		pointOwnership[i] = make(map[int]int)
		for j := 0; j < maxY; j++ {
			closestId := closestPoint(allPoints, i, j)
			if closestId == -1 {
				pointOwnership[i][j] = -1
				continue
			}
			pointOwnership[i][j] = closestId
			closestHisto[closestId] += 1
		}
	}

	for i := 0; i < maxX; i++ {
		topOwner := pointOwnership[i][0]
		bottomOwner := pointOwnership[i][maxY]
		delete(closestHisto, topOwner)
		delete(closestHisto, bottomOwner)
	}

	for i := 0; i < maxY; i++ {
		leftOwner := pointOwnership[0][i]
		rightOwner := pointOwnership[maxX][i]
		delete(closestHisto, leftOwner)
		delete(closestHisto, rightOwner)
	}

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
