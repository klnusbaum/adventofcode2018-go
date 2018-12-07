package main

/*
As you're about to begin construction, four of the Elves offer to help. "The sun will set soon; it'll go faster if we work together." Now, you need to account for multiple people working on steps simultaneously. If multiple steps are available, workers should still begin them in alphabetical order.

Each step takes 60 seconds plus an amount corresponding to its letter: A=1, B=2, C=3, and so on. So, step A takes 60+1=61 seconds, while step Z takes 60+26=86 seconds. No time is required between steps.

To simplify things for the example, however, suppose you only have help from one Elf (a total of two workers) and that each step takes 60 fewer seconds (so that step A takes 1 second and step Z takes 26 seconds). Then, using the same instructions as above, this is how each second would be spent:

Second   Worker 1   Worker 2   Done
   0        C          .
   1        C          .
   2        C          .
   3        A          F       C
   4        B          F       CA
   5        B          F       CA
   6        D          F       CAB
   7        D          F       CAB
   8        D          F       CAB
   9        D          .       CABF
  10        E          .       CABFD
  11        E          .       CABFD
  12        E          .       CABFD
  13        E          .       CABFD
  14        E          .       CABFD
  15        .          .       CABFDE
Each row represents one second of time. The Second column identifies how many seconds have passed as of the beginning of that second. Each worker column shows the step that worker is currently doing (or . if they are idle). The Done column shows completed steps.

Note that the order of the steps has changed; this is because steps now take time to finish and multiple workers can begin multiple steps simultaneously.

In this example, it would take 15 seconds for two workers to complete these steps.

With 5 workers and the 60+ second step durations described above, how long will it take to complete all of the steps?


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
	for {
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

		sub1(workerTimeLeft)
		if len(requirements) == 0 && allZero(workerTimeLeft) {
			break
		}
		t++
	}

	return t
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
