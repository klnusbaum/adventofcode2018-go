package main

/*
--- Day 5: Alchemical Reduction ---
You've managed to sneak in to the prototype suit manufacturing lab. The Elves are making decent progress, but are still struggling with the suit's size reduction capabilities.

While the very latest in 1518 alchemical technology might have solved their problem eventually, you can do better. You scan the chemical composition of the suit's material and discover that it is formed by extremely long polymers (one of which is available as your puzzle input).

The polymer is formed by smaller units which, when triggered, react with each other such that two adjacent units of the same type and opposite polarity are destroyed. Units' types are represented by letters; units' polarity is represented by capitalization. For instance, r and R are units with the same type but opposite polarity, whereas r and s are entirely different types and do not react.

For example:

In aA, a and A react, leaving nothing behind.
In abBA, bB destroys itself, leaving aA. As above, this then destroys itself, leaving nothing.
In abAB, no two adjacent units are of the same type, and so nothing happens.
In aabAAB, even though aa and AA are of the same type, their polarities match, and so nothing happens.
Now, consider a larger example, dabAcCaCBAcCcaDA:

dabAcCaCBAcCcaDA  The first 'cC' is removed.
dabAaCBAcCcaDA    This creates 'Aa', which is removed.
dabCBAcCcaDA      Either 'cC' or 'Cc' are removed (the result is the same).
dabCBAcaDA        No further actions can be taken.
After all possible reactions, the resulting polymer contains 10 units.

How many units remain after fully reacting the polymer you scanned? (Note: in this puzzle and others, the input is large; if you copy/paste your input, make sure you get the whole thing.)
*/

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

func main() {
	filename := os.Args[1]

	if err := react(filename); err != nil {
		fmt.Printf("Error analyzing schedules: %v\n", err)
		os.Exit(1)
	}
}

func react(filename string) error {
	polymer, err := readPolymer(filename)
	if err != nil {
		return fmt.Errorf("couldn't read polymer from file %s: %v", filename, err)
	}
	processed := processPolymer(polymer)
	fmt.Printf("Number of units: %d\n", len(processed))
	return nil
}

func readPolymer(filename string) (string, error) {
	file, _ := os.Open(filename)
	contents, _ := ioutil.ReadAll(file)
	return strings.TrimSpace(string(contents)), nil
}

// This could vastly be impoved by instead of remove characters,
// just building up a new string for each pass over the polymer
func processPolymer(polymer string) string {
	toReprocess := polymer
	for i := 0; i+1 < len(polymer); i++ {
		if cancelOut(polymer[i], polymer[i+1]) {
			toReprocess = polymer[0:i] + polymer[i+2:]
			break
		}
	}

	if toReprocess == polymer {
		return toReprocess
	}
	return processPolymer(toReprocess)
}

func cancelOut(b1, b2 byte) bool {
	r1, r2 := rune(b1), rune(b2)
	return (unicode.IsLower(r1) && unicode.ToUpper(r1) == r2) ||
		(unicode.IsUpper(r1) && unicode.ToLower(r1) == r2)
}
