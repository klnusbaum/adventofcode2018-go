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

var alphaTimes = map[string]int{
	"A": 61,
	"B": 62,
	"C": 63,
	"D": 64,
	"E": 65,
	"F": 66,
	"G": 67,
	"H": 68,
	"I": 69,
	"J": 70,
	"K": 71,
	"L": 72,
	"M": 73,
	"N": 74,
	"O": 75,
	"P": 76,
	"Q": 77,
	"R": 78,
	"S": 79,
	"T": 80,
	"U": 81,
	"V": 82,
	"W": 83,
	"X": 84,
	"Y": 85,
	"Z": 86,
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
		return fmt.Errorf("Error loading lines: %v", lines)
	}
	requirements := loadReqs(lines)

	totalTime := doSteps(requirements)

	fmt.Printf("Total time: %d\n", totalTime)

	return nil
}

func doSteps(requirements map[string]map[string]bool) int {
	t := 0
	workerLetters := []string{"", "", "", "", ""}
	workerTimeLeft := []int{0, 0, 0, 0, 0}
	for !(len(requirements) == 0 && allZero(workerTimeLeft)) {
		for i, timeLeft := range workerTimeLeft {
			if timeLeft == 0 {
				finishedLetter := workerLetters[i]

				for _, req := range requirements {
					delete(req, finishedLetter)
				}

				workerLetters[i] = ""
				delete(requirements, finishedLetter)
			}
		}

		var readyLetters []string
		for letter, reqs := range requirements {
			if len(reqs) == 0 {
				readyLetters = append(readyLetters, letter)
			}
		}

		sort.Strings(readyLetters)
		var lettersNotBeingWorkedOn []string
		for _, letter := range readyLetters {
			if letter != workerLetters[0] &&
				letter != workerLetters[1] &&
				letter != workerLetters[2] &&
				letter != workerLetters[3] &&
				letter != workerLetters[4] {
				lettersNotBeingWorkedOn = append(lettersNotBeingWorkedOn, letter)
			}
		}

		for _, letter := range lettersNotBeingWorkedOn {
			for i, _ := range workerLetters {
				if workerLetters[i] == "" {
					workerLetters[i] = letter
					workerTimeLeft[i] = alphaTimes[letter]
					break
				}
			}
		}

		t++
		sub1(workerTimeLeft)
	}

	return t - 1
}

func sub1(times []int) {
	for i, time := range times {
		if time > 0 {
			times[i] -= 1
		}
	}
}

func allZero(times []int) bool {
	for _, time := range times {
		if time > 0 {
			return false
		}
	}
	return true
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
