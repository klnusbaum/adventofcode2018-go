package main

/*
Time to improve the polymer.

One of the unit types is causing problems; it's preventing the polymer from collapsing as much as it should. Your goal is to figure out which unit type is causing the most problems, remove all instances of it (regardless of polarity), fully react the remaining polymer, and measure its length.

For example, again using the polymer dabAcCaCBAcCcaDA from above:

Removing all A/a units produces dbcCCBcCcD. Fully reacting this polymer produces dbCBcD, which has length 6.
Removing all B/b units produces daAcCaCAcCcaDA. Fully reacting this polymer produces daCAcaDA, which has length 8.
Removing all C/c units produces dabAaBAaDA. Fully reacting this polymer produces daDA, which has length 4.
Removing all D/d units produces abAcCaCBAcCcaA. Fully reacting this polymer produces abCBAc, which has length 6.
In this example, removing all C/c units was best, producing the answer 4.

What is the length of the shortest polymer you can produce by removing all units of exactly one type and fully reacting the result?
*/

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

var lowerAlpha = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
var upperAlpha = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func main() {
	if err := react("input.txt"); err != nil {
		fmt.Printf("Error reacting: %v\n", err)
		os.Exit(1)
	}
}

func react(filename string) error {
	polymer, err := readPolymer(filename)
	if err != nil {
		return fmt.Errorf("couldn't read polymer from file %s: %v", filename, err)
	}
	permutedPolymers := permutePolymer(polymer)
	shortest := shortestLength(permutedPolymers)
	fmt.Printf("Length of shortest polymer: %d\n", shortest)
	return nil
}

func readPolymer(filename string) (string, error) {
	file, _ := os.Open(filename)
	contents, _ := ioutil.ReadAll(file)
	return strings.TrimSpace(string(contents)), nil
}

func permutePolymer(polymer string) []string {
	var polymers []string
	for i, letter := range lowerAlpha {
		toProcess := strings.Replace(polymer, letter, "", -1)
		toProcess = strings.Replace(toProcess, upperAlpha[i], "", -1)
		polymers = append(polymers, processPolymer(toProcess))
	}
	return polymers
}

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

func shortestLength(polymers []string) int {
	shortLength := 100000000000
	for _, polymer := range polymers {
		if len(polymer) < shortLength {
			shortLength = len(polymer)
		}
	}
	return shortLength
}
