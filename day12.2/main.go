package main

// Start Dec 13 10:35 am
// Part 1 Dec 13 11:31 am
/* Part 2 Dec 13 14:27
 */

import (
	"fmt"
	"github.com/klnusbaum/adventofcode2018/ll"
	"os"
	"strings"
	"unicode/utf8"
)

type mask struct {
	ll bool
	l  bool
	c  bool
	r  bool
	rr bool
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

	state, initialLength, rules := loadStateAndRules(lines)
	minBound := -2
	maxBound := initialLength + 2
	g := 1
	for ; g <= 240; g++ {
		oldState := state
		state = calcNewState(state, rules, minBound, maxBound)
		if state[minBound] || state[minBound+1] {
			minBound -= 2
		} else if !(state[minBound] && state[minBound+1] && state[minBound+2]) {
			minBound += 1
		}
		if state[maxBound-2] || state[maxBound-1] || state[maxBound] || state[maxBound+1] {
			maxBound += 2
		}

		if inSteadyState(oldState, state) {
			break
		}
	}

	finalInts := make([]int, 0)
	for potNum, state := range state {
		if state {
			finalInts = append(finalInts, potNum+(50000000000)-g)
		}
	}

	fmt.Printf("Sum of pot numbers is: %d\n", sumInts(finalInts))

	return nil
}

func potSum(state map[int]bool) int {
	sum := 0
	for potNum, hasPlant := range state {
		if hasPlant {
			sum += potNum
		}
	}
	return sum
}

func calcNewState(state map[int]bool, rules map[mask]bool, minBound, maxBound int) map[int]bool {
	newState := make(map[int]bool)
	for i := minBound; i < maxBound; i++ {
		env := mask{
			ll: state[i-2],
			l:  state[i-1],
			c:  state[i],
			r:  state[i+1],
			rr: state[i+2],
		}
		newState[i] = rules[env]
	}

	return newState
}

func loadStateAndRules(lines []string) (map[int]bool, int, map[mask]bool) {
	init, maxPos := loadState(lines[0])
	rules := loadRules(lines[2:])

	return init, maxPos, rules
}

func loadState(toLoad string) (map[int]bool, int) {
	toLoad = strings.TrimPrefix(toLoad, "initial state: ")
	state := make(map[int]bool)
	for i, char := range toLoad {
		state[i] = dotPoundtoBool(char)
	}
	return state, len(toLoad)
}

func loadRules(lines []string) map[mask]bool {
	rules := make(map[mask]bool)
	for _, line := range lines {
		m, res := loadRule(line)
		rules[m] = res
	}
	return rules
}

func loadRule(line string) (mask, bool) {
	components := strings.Split(line, " => ")
	bools := make([]bool, 5)
	for i, char := range components[0] {
		bools[i] = dotPoundtoBool(char)
	}

	m := mask{
		ll: bools[0],
		l:  bools[1],
		c:  bools[2],
		r:  bools[3],
		rr: bools[4],
	}

	firstRune, _ := utf8.DecodeRuneInString(components[1])
	result := dotPoundtoBool(firstRune)

	return m, result

}

func dotPoundtoBool(r rune) bool {
	if r == '#' {
		return true
	}
	return false
}

func printState(state map[int]bool, min, max, gen int) {
	for i := min; i < max; i++ {
		if state[i] {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}
	}
	fmt.Printf("%d, %d, %d\n", min, max, gen)
}

func inSteadyState(state, newState map[int]bool) bool {
	for potNum, hasFlower := range state {
		if hasFlower != newState[potNum+1] {
			return false
		}
	}

	return true
}

func sumInts(nums []int) int {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	return sum
}
