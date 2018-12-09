package main

/*
--- Part Two ---
Amused by the speed of your answer, the Elves are curious:

What would the new winning Elf's score be if the number of the last marble were 100 times larger?
*/

import (
	"fmt"
)

func main() {
	// analyze(459, 72103)
	analyze(459, 72103*100)
	// analyze(9, 25)
	// analyze(10, 1618)
}

type node struct {
	next *node
	prev *node
	val  int
}

func printNodes(n *node) {
	for n.val != 0 {
		n = n.next
	}

	fmt.Printf("[ %d, ", n.val)
	n = n.next
	for n.val != 0 {
		fmt.Printf("%d, ", n.val)
		n = n.next
	}

	fmt.Printf("]\n")
}

func analyze(numPlayers, numMarbles int) {
	elfScores := make([]int, numPlayers)
	currentElf := 0
	currentMarble := &node{
		val: 0,
	}
	currentMarble.next = currentMarble
	currentMarble.prev = currentMarble
	// printNodes(currentMarble)
	for i := 1; i <= numMarbles; i++ {
		currentElf++
		if currentElf == numPlayers {
			currentElf = 0
		}

		if i%23 == 0 {
			elfScores[currentElf] += i
			toRemove := currentMarble.prev.prev.prev.prev.prev.prev.prev
			elfScores[currentElf] += toRemove.val

			theNext := toRemove.next
			thePrev := toRemove.prev

			thePrev.next = theNext
			theNext.prev = thePrev

			currentMarble = theNext
			continue
		}

		newNext := currentMarble.next.next
		newPrev := currentMarble.next
		newMarble := &node{
			next: newNext,
			prev: newPrev,
			val:  i,
		}

		newNext.prev = newMarble
		newPrev.next = newMarble
		currentMarble = newMarble
		// printNodes(currentMarble)
	}

	fmt.Printf("Max score: %d\n", maxScore(elfScores))
}

func maxScore(scores []int) int {
	ms := scores[0]
	for i := 1; i < len(scores); i++ {
		if scores[i] > ms {
			ms = scores[i]
		}
	}
	return ms
}
