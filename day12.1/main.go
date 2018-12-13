package main

// start Dec 13 10:35 am
// finish part 1  Dev 13 11:31 am
/*
 */

import (
	"fmt"
	"github.com/klnusbaum/adventofcode2018/ll"
	"os"
	"strings"
	"unicode/utf8"
)

const (
	_negPadding = 2
)

type rule struct {
	mask   []bool
	result bool
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
	for i := 1; i < 21; i++ {
		newState := make(map[int]bool)
		for i := minBound; i < maxBound; i++ {
			env := []bool{state[i-2], state[i-1], state[i], state[i+1], state[i+2]}
			matchingRule := findMatchingRule(env, rules)
			newState[i] = matchingRule.result
		}
		state = newState
		minBound = minBound - 2
		maxBound = maxBound + 2
	}

	sum := 0
	for potNum, hasPlant := range state {
		if hasPlant {
			sum += potNum
		}
	}

	fmt.Printf("Sum of pot numbers is: %d\n", sum)

	return nil
}

func loadStateAndRules(lines []string) (map[int]bool, int, []rule) {
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

func loadRules(lines []string) []rule {
	rules := make([]rule, len(lines))
	for i, line := range lines {
		rules[i] = loadRule(line)
	}
	return rules
}

func loadRule(line string) rule {
	components := strings.Split(line, " => ")
	mask := make([]bool, 5)
	for i, char := range components[0] {
		mask[i] = dotPoundtoBool(char)
	}

	firstRune, _ := utf8.DecodeRuneInString(components[1])
	result := dotPoundtoBool(firstRune)

	return rule{
		mask:   mask,
		result: result,
	}

}

func dotPoundtoBool(r rune) bool {
	if r == '#' {
		return true
	}
	return false
}

func findMatchingRule(env []bool, rules []rule) rule {
	for _, rule := range rules {
		if boolArraysMatch(env, rule.mask) {
			return rule
		}
	}
	panic("Couldn't find matching rule")
}

func boolArraysMatch(one, two []bool) bool {
	for i, b1 := range one {
		if b1 != two[i] {
			return false
		}
	}
	return true
}
