package main

// start Dec 14 14:12
// finish 14:55
/* finish 15:19
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
	inputSeq := getSeq(input)
	elfOneIndex := 0
	elfTwoIndex := 1
	scoreboard := []int{3, 7}
	scoreSeqPos := -1
	for scoreSeqPos == -1 {
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

		scoreSeqPos = getScoreSecPos(scoreboard, inputSeq)

		// fmt.Printf("Scoreboard: %v\n", scoreboard)
		// fmt.Printf("Indexes: %d %d\n", elfOneIndex, elfTwoIndex)
	}

	fmt.Printf("Num to the left: %d\n", scoreSeqPos)

	return nil
}

func getScoreSecPos(scoreboard []int, sequence []int) int {
	if len(sequence) > len(scoreboard) {
		return -1
	}

	lastSeq := scoreboard[len(scoreboard)-len(sequence):]
	// fmt.Printf("Scoreboard: %v\n", scoreboard)
	// fmt.Printf("Comparing: %v and %v\n", lastSeq, sequence)
	if slicesEqual(lastSeq, sequence) {
		return len(scoreboard) - len(sequence)
	}

	if len(sequence)+1 > len(scoreboard) {
		return -1
	}

	secondToLastSeq := scoreboard[len(scoreboard)-len(sequence)-1 : len(scoreboard)-1]
	// fmt.Printf("Comparing: %v and %v\n", secondToLastSeq, sequence)
	if slicesEqual(secondToLastSeq, sequence) {
		return len(scoreboard) - len(sequence) - 1
	}

	return -1
}

func getSeq(num int) []int {
	var seq []int
	for num > 0 {
		digit := num % 10
		seq = append([]int{digit}, seq...)
		num /= 10
	}

	return seq
}

func slicesEqual(s1, s2 []int) bool {
	for i, elem := range s1 {
		if elem != s2[i] {
			return false
		}

	}
	return true
}
