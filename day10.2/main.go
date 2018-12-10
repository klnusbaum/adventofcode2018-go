package main

/*
--- Part Two ---
Good thing you didn't have to wait, because that would have taken a long time - much longer than the 3 seconds in the example above.

Impressed by your sub-hour communication capabilities, the Elves are curious: exactly how many seconds would they have needed to wait for that message to appear?
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
