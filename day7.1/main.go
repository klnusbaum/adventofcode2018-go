package main

/*
--- Day 7: The Sum of Its Parts ---
You find yourself standing on a snow-covered coastline; apparently, you landed a little off course. The region is too hilly to see the North Pole from here, but you do spot some Elves that seem to be trying to unpack something that washed ashore. It's quite cold out, so you decide to risk creating a paradox by asking them for directions.

"Oh, are you the search party?" Somehow, you can understand whatever Elves from the year 1018 speak; you assume it's Ancient Nordic Elvish. Could the device on your wrist also be a translator? "Those clothes don't look very warm; take this." They hand you a heavy coat.

"We do need to find our way back to the North Pole, but we have higher priorities at the moment. You see, believe it or not, this box contains something that will solve all of Santa's transportation problems - at least, that's what it looks like from the pictures in the instructions." It doesn't seem like they can read whatever language it's in, but you can: "Sleigh kit. Some assembly required."

"'Sleigh'? What a wonderful name! You must help us assemble this 'sleigh' at once!" They start excitedly pulling more parts out of the box.

The instructions specify a series of steps and requirements about which steps must be finished before others can begin (your puzzle input). Each step is designated by a single letter. For example, suppose you have the following instructions:

Step C must be finished before step A can begin.
Step C must be finished before step F can begin.
Step A must be finished before step B can begin.
Step A must be finished before step D can begin.
Step B must be finished before step E can begin.
Step D must be finished before step E can begin.
Step F must be finished before step E can begin.
Visually, these requirements look like this:


  -->A--->B--
 /    \      \
C      -->D----->E
 \           /
  ---->F-----
Your first goal is to determine the order in which the steps should be completed. If more than one step is ready, choose the step which is first alphabetically. In this example, the steps would be completed as follows:

Only C is available, and so it is done first.
Next, both A and F are available. A is first alphabetically, so it is done next.
Then, even though F was available earlier, steps B and D are now also available, and B is the first alphabetically of the three.
After that, only D and F are available. E is not available because only some of its prerequisites are complete. Therefore, D is completed next.
F is the only choice, so it is done next.
Finally, E is completed.
So, in this example, the correct order is CABDFE.

In what order should the steps in your instructions be completed?

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
