package main

/*
You discover a dial on the side of the device; it seems to let you select a square of any size, not just 3x3. Sizes from 1x1 to 300x300 are supported.

Realizing this, you now must find the square of any size with the largest total power. Identify this square by including its size as a third parameter after the top-left coordinate: a 9x9 square with a top-left corner of 3,5 is identified as 3,5,9.

For example:

For grid serial number 18, the largest total square (with a total power of 113) is 16x16 and has a top-left corner of 90,269, so its identifier is 90,269,16.
For grid serial number 42, the largest total square (with a total power of 119) is 12x12 and has a top-left corner of 232,251, so its identifier is 232,251,12.
What is the X,Y,size identifier of the square with the largest total power?
*/

import (
	"fmt"
	"os"
)

func main() {
	serialNumber := 8868
	if err := analyze(serialNumber); err != nil {
		fmt.Printf("Error analyzing: %v\n", err)
		os.Exit(1)
	}
}

func analyze(serialNumber int) error {
	powerLevels := getPowerLevels(serialNumber)
	x, y, size := largestPowerSquare(powerLevels)
	fmt.Printf("Most power grid is: %d,%d,%d\n", x, y, size)

	return nil
}

func getPowerLevels(serialNumber int) [][]int {
	powerLevels := make([][]int, 300)
	for i, _ := range powerLevels {
		powerLevels[i] = make([]int, 300)
	}

	for i := 0; i < 300; i++ {
		for j := 0; j < 300; j++ {
			powerLevels[i][j] = getPowerLevel(i+1, j+1, serialNumber)
		}
	}

	return powerLevels
}

func getPowerLevel(x, y, serialNumber int) int {
	rackID := x + 10
	powerLevel := rackID * y
	powerLevel += serialNumber
	powerLevel *= rackID
	powerLevel = (powerLevel / 100) % 10
	return powerLevel - 5
}

func largestPowerSquare(powerLevels [][]int) (int, int, int) {
	largestX := -1
	largestY := -1
	largestPower := 0
	largestSquareSize := 0
	for squareSize := 1; squareSize <= 300; squareSize++ {
		for i := 0; i < 301-squareSize; i++ {
			for j := 0; j < 301-squareSize; j++ {
				power := getSquarePower(powerLevels, i, j, squareSize)
				if power > largestPower {
					largestPower = power
					largestX = i + 1
					largestY = j + 1
					largestSquareSize = squareSize
				}
			}
		}
	}

	return largestX, largestY, largestSquareSize
}

func getSquarePower(powerLevels [][]int, i, j, squareSize int) int {
	power := 0
	for ii := i; ii < i+squareSize; ii++ {
		for jj := j; jj < j+squareSize; jj++ {
			power += powerLevels[ii][jj]
		}
	}

	return power
}
