package main

/*
 */

import (
	"fmt"
	"github.com/klnusbaum/adventofcode2018/ll"
	"os"
	"sort"
	"strings"
)

var upperAlpha = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

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
		return fmt.Errorf("Error loading lines: %v", lines)
	}
	requirements := loadReqs(lines)
	// for step, reqs := range requirements {
	// 	fmt.Printf("Step %s requires %v\n", step, reqs)
	// }

	order := getReadySteps(requirements)

	fmt.Printf("Order is: %s\n", strings.Join(order, ""))

	return nil
}

func getReadySteps(requirements map[string]map[string]bool) []string {
	if len(requirements) == 0 {
		return []string{}
	}

	var readyLetters []string
	for letter, reqs := range requirements {
		if len(reqs) == 0 {
			readyLetters = append(readyLetters, letter)
		}
	}

	sort.Strings(readyLetters)
	next := readyLetters[0]
	fmt.Printf("Next: %s\n", next)

	for letter, _ := range requirements {
		delete(requirements[letter], next)
	}

	delete(requirements, next)

	return append([]string{next}, getReadySteps(requirements)...)
}

func loadReqs(lines []string) map[string]map[string]bool {
	reqs := make(map[string]map[string]bool)
	for _, letter := range upperAlpha {
		reqs[letter] = make(map[string]bool)
	}

	for _, line := range lines {
		fields := strings.Fields(line)
		reqs[fields[7]][fields[1]] = true
	}

	return reqs
}
