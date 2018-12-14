package main

// start Dec 14 14:12
// finish 14:55
/*
 */

import (
	"fmt"
	"os"
)

func main() {
	if err := analyze(540391); err != nil {
		fmt.Printf("Error analyzing: %v\n", err)
		os.Exit(1)
	}
}

func analyze(input int) error {
	elfOneIndex := 0
	elfTwoIndex := 1
	scoreboard := []int{3, 7}
	for len(scoreboard) < input+10 {
		newScore := scoreboard[elfOneIndex] + scoreboard[elfTwoIndex]
		if newScore > 9 {
			first := newScore / 10
			second := newScore % 10
			scoreboard = append(scoreboard, first, second)
		} else {
			scoreboard = append(scoreboard, newScore)
		}

		elfOneSteps := 1 + scoreboard[elfOneIndex]
		elfTwoSteps := 1 + scoreboard[elfTwoIndex]
		// fmt.Printf("Steps: %d, %d\n", elfOneSteps, elfTwoSteps)
		for elfOneSteps > 0 {
			if elfOneIndex == len(scoreboard)-1 {
				elfOneIndex = 0
			} else {
				elfOneIndex++
			}
			elfOneSteps--
		}
		for elfTwoSteps > 0 {
			if elfTwoIndex == len(scoreboard)-1 {
				elfTwoIndex = 0
			} else {
				elfTwoIndex++
			}
			elfTwoSteps--
		}

		// fmt.Printf("Scoreboard: %v\n", scoreboard)
		// fmt.Printf("Indexes: %d %d\n", elfOneIndex, elfTwoIndex)
	}

	last10elems := scoreboard[len(scoreboard)-10:]
	for _, elem := range last10elems {
		fmt.Printf("%d", elem)
	}
	fmt.Println()

	return nil
}
