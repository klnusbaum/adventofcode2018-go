package main

/*
--- Day 10: The Stars Align ---
It's no use; your navigation system simply isn't capable of providing walking directions in the arctic circle, and certainly not in 1018.

The Elves suggest an alternative. In times like these, North Pole rescue operations will arrange points of light in the sky to guide missing Elves back to base. Unfortunately, the message is easy to miss: the points move slowly enough that it takes hours to align them, but have so much momentum that they only stay aligned for a second. If you blink at the wrong time, it might be hours before another message appears.

You can see these points of light floating in the distance, and record their position in the sky and their velocity, the relative change in position per second (your puzzle input). The coordinates are all given from your perspective; given enough time, those positions and velocities will move the points into a cohesive message!

Rather than wait, you decide to fast-forward the process and calculate what the points will eventually spell.

For example, suppose you note the following points:

position=< 9,  1> velocity=< 0,  2>
position=< 7,  0> velocity=<-1,  0>
position=< 3, -2> velocity=<-1,  1>
position=< 6, 10> velocity=<-2, -1>
position=< 2, -4> velocity=< 2,  2>
position=<-6, 10> velocity=< 2, -2>
position=< 1,  8> velocity=< 1, -1>
position=< 1,  7> velocity=< 1,  0>
position=<-3, 11> velocity=< 1, -2>
position=< 7,  6> velocity=<-1, -1>
position=<-2,  3> velocity=< 1,  0>
position=<-4,  3> velocity=< 2,  0>
position=<10, -3> velocity=<-1,  1>
position=< 5, 11> velocity=< 1, -2>
position=< 4,  7> velocity=< 0, -1>
position=< 8, -2> velocity=< 0,  1>
position=<15,  0> velocity=<-2,  0>
position=< 1,  6> velocity=< 1,  0>
position=< 8,  9> velocity=< 0, -1>
position=< 3,  3> velocity=<-1,  1>
position=< 0,  5> velocity=< 0, -1>
position=<-2,  2> velocity=< 2,  0>
position=< 5, -2> velocity=< 1,  2>
position=< 1,  4> velocity=< 2,  1>
position=<-2,  7> velocity=< 2, -2>
position=< 3,  6> velocity=<-1, -1>
position=< 5,  0> velocity=< 1,  0>
position=<-6,  0> velocity=< 2,  0>
position=< 5,  9> velocity=< 1, -2>
position=<14,  7> velocity=<-2,  0>
position=<-3,  6> velocity=< 2, -1>
Each line represents one point. Positions are given as <X, Y> pairs: X represents how far left (negative) or right (positive) the point appears, while Y represents how far up (negative) or down (positive) the point appears.

At 0 seconds, each point has the position given. Each second, each point's velocity is added to its position. So, a point with velocity <1, -2> is moving to the right, but is moving upward twice as quickly. If this point's initial position were <3, 9>, after 3 seconds, its position would become <6, 3>.

Over time, the points listed above would move like this:

Initially:
........#.............
................#.....
.........#.#..#.......
......................
#..........#.#.......#
...............#......
....#.................
..#.#....#............
.......#..............
......#...............
...#...#.#...#........
....#..#..#.........#.
.......#..............
...........#..#.......
#...........#.........
...#.......#..........

After 1 second:
......................
......................
..........#....#......
........#.....#.......
..#.........#......#..
......................
......#...............
....##.........#......
......#.#.............
.....##.##..#.........
........#.#...........
........#...#.....#...
..#...........#.......
....#.....#.#.........
......................
......................

After 2 seconds:
......................
......................
......................
..............#.......
....#..#...####..#....
......................
........#....#........
......#.#.............
.......#...#..........
.......#..#..#.#......
....#....#.#..........
.....#...#...##.#.....
........#.............
......................
......................
......................

After 3 seconds:
......................
......................
......................
......................
......#...#..###......
......#...#...#.......
......#...#...#.......
......#####...#.......
......#...#...#.......
......#...#...#.......
......#...#...#.......
......#...#..###......
......................
......................
......................
......................

After 4 seconds:
......................
......................
......................
............#.........
........##...#.#......
......#.....#..#......
.....#..##.##.#.......
.......##.#....#......
...........#....#.....
..............#.......
....#......#...#......
.....#.....##.........
...............#......
...............#......
......................
......................
After 3 seconds, the message appeared briefly: HI. Of course, your message will be much longer and will take many more seconds to appear.

What message will eventually appear in the sky?
*/

import (
	"fmt"
	"github.com/klnusbaum/adventofcode2018/ll"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"strings"
)

type point struct {
	posX int
	posY int
	velX int
	velY int
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

	points := loadPoints(lines)
	numIters := 1
	for !withinBounds(points) {
		for _, point := range points {
			point.posX += point.velX
			point.posY += point.velY
		}

		numIters++
	}

	for withinBounds(points) {
		for _, point := range points {
			point.posX += point.velX
			point.posY += point.velY
		}

		printImage(points, numIters)

		numIters++
	}

	return nil
}

func loadPoints(lines []string) []*point {
	points := make([]*point, len(lines))
	for i, line := range lines {
		pos := line[:25]
		vel := line[26:]

		pos = strings.TrimPrefix(pos, "position=<")
		pos = strings.TrimSpace(strings.TrimSuffix(pos, ">"))
		xandy := strings.Split(pos, ",")
		x, _ := strconv.Atoi(strings.TrimSpace(xandy[0]))
		y, _ := strconv.Atoi(strings.TrimSpace(xandy[1]))

		vel = strings.TrimPrefix(vel, "velocity=<")
		vel = strings.TrimSuffix(vel, ">")
		velxandy := strings.Split(vel, ",")
		velx, _ := strconv.Atoi(strings.TrimSpace(velxandy[0]))
		vely, _ := strconv.Atoi(strings.TrimSpace(velxandy[1]))
		points[i] = &point{
			posX: x,
			posY: y,
			velX: velx,
			velY: vely,
		}
	}

	return points
}

func withinBounds(points []*point) bool {
	maxX, minX, maxY, minY := points[0].posX, points[0].posX, points[0].posY, points[0].posY
	for i := 1; i < len(points); i++ {
		if points[i].posX > maxX {
			maxX = points[i].posX
		} else if points[i].posX < minX {
			minX = points[i].posX
		}

		if points[i].posY > maxY {
			maxY = points[i].posY
		} else if points[i].posY < minY {
			minY = points[i].posY
		}
	}

	return abs(maxX-minX) < 100 && abs(maxY-minY) < 100
}

func printImage(points []*point, i int) {
	image := image.NewRGBA(
		image.Rectangle{
			Min: image.Point{
				X: 0,
				Y: 0,
			},
			Max: image.Point{
				X: 300,
				Y: 200,
			},
		},
	)

	for i := 0; i < 300; i++ {
		for j := 0; j < 300; j++ {
			image.Set(i, j, color.RGBA{0, 0, 0, 255})
		}
	}

	for _, point := range points {
		image.Set(point.posX, point.posY, color.RGBA{255, 100, 0, 255})
	}

	f, _ := os.Create(fmt.Sprintf("output/image%d.png", i))
	png.Encode(f, image)
	f.Close()
}

func abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}
