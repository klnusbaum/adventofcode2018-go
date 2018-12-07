package main

/*
On the other hand, if the coordinates are safe, maybe the best you can do is try to find a region near as many coordinates as possible.

For example, suppose you want the sum of the Manhattan distance to all of the coordinates to be less than 32. For each location, add up the distances to all of the given coordinates; if the total of those distances is less than 32, that location is within the desired region. Using the same coordinates as above, the resulting region looks like this:

..........
.A........
..........
...###..C.
..#D###...
..###E#...
.B.###....
..........
..........
........F.
In particular, consider the highlighted location 4,3 located at the top middle of the region. Its calculation is as follows, where abs() is the absolute value function:

Distance to coordinate A: abs(4-1) + abs(3-1) =  5
Distance to coordinate B: abs(4-1) + abs(3-6) =  6
Distance to coordinate C: abs(4-8) + abs(3-3) =  4
Distance to coordinate D: abs(4-3) + abs(3-4) =  2
Distance to coordinate E: abs(4-5) + abs(3-5) =  3
Distance to coordinate F: abs(4-8) + abs(3-9) = 10
Total distance: 5 + 6 + 4 + 2 + 3 + 10 = 30
Because the total distance to all coordinates (30) is less than 32, the location is within the region.

This region, which also includes coordinates D and E, has a total size of 16.

Your actual region will need to be much larger than this example, though, instead including all locations with a total distance of less than 10000.

What is the size of the region containing all locations which have a total distance to all given coordinates of less than 10000?
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

	area := 0

	for i := 0; i < maxX; i++ {
		for j := 0; j < maxY; j++ {
			totalDistanceFromAllPoints := 0
			for _, point := range allPoints {
				totalDistanceFromAllPoints += manDistance(i, j, point.x, point.y)
			}

			if totalDistanceFromAllPoints < 10000 {
				area++
			}
		}
	}

	fmt.Printf("The area is: %d\n", area)
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

func manDistance(x1, y1, x2, y2 int) int {
	return abs(x1-x2) + abs(y1-y2)
}

func abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
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
