/*
You notice that the device repeats the same frequency change list over and over. To calibrate the device, you need to find the first frequency it reaches twice.

For example, using the same list of changes above, the device would loop as follows:

Current frequency  0, change of +1; resulting frequency  1.
Current frequency  1, change of -2; resulting frequency -1.
Current frequency -1, change of +3; resulting frequency  2.
Current frequency  2, change of +1; resulting frequency  3.
(At this point, the device continues from the start of the list.)
Current frequency  3, change of +1; resulting frequency  4.
Current frequency  4, change of -2; resulting frequency  2, which has already been seen.
In this example, the first frequency reached twice is 2. Note that your device might need to repeat its list of frequency changes many times before a duplicate frequency is found, and that duplicates might be found while in the middle of processing the list.

Here are other examples:

+1, -1 first reaches 0 twice.
+3, +3, +4, -2, -4 first reaches 10 twice.
-6, +3, +8, +5, -6 first reaches 5 twice.
+7, +7, -2, -7, -4 first reaches 14 twice.
What is the first frequency your device reaches twice?
*/

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	filename := os.Args[1]
	acc := newAccumulator(filename)

	if err := acc.accumulate(os.Stdout); err != nil {
		fmt.Printf("Couldn't accumulate result: %v", err)
		os.Exit(1)
	}
}

func newAccumulator(filename string) *accumulator {
	return &accumulator{
		filename:  filename,
		lines:     make([]string, 0),
		seenFreqs: make(map[int]struct{}),
		counter:   0,
	}
}

type accumulator struct {
	filename  string
	lines     []string
	seenFreqs map[int]struct{}
	counter   int
}

func (acc *accumulator) accumulate(out io.Writer) error {
	acc.loadLines()

	var err error
	hasSeen := false
	i := 0
	for !hasSeen {
		hasSeen, err = acc.checkFreq(acc.lines[i])
		if err != nil {
			return fmt.Errorf("Couldn't check frequency at line %q: %v", i, err)
		}

		i++
		if i >= len(acc.lines) {
			i = 0
		}
	}

	_, err = fmt.Fprintf(out, "Repeated frequency: %d\n", acc.counter)
	return err
}

func (acc *accumulator) checkFreq(line string) (bool, error) {
	if err := acc.accLine(line); err != nil {
		return false, err
	}

	_, hasSeen := acc.seenFreqs[acc.counter]
	if hasSeen {
		return true, nil
	}

	acc.seenFreqs[acc.counter] = struct{}{}
	return false, nil
}

func (acc *accumulator) accLine(line string) error {
	switch firstChar := string(line[0]); firstChar {
	case "+":
		return acc.addNum(line[1:])
	case "-":
		return acc.subNum(line[1:])
	default:
		return fmt.Errorf("Unrecognized first character %q on line %q", firstChar, line)
	}
}

func (acc *accumulator) addNum(num string) error {
	toAdd, err := strconv.Atoi(num)
	if err != nil {
		return acc.numConvErr(num)
	}
	acc.counter += toAdd
	return nil
}

func (acc *accumulator) subNum(num string) error {
	toSub, err := strconv.Atoi(num)
	if err != nil {
		return acc.numConvErr(num)
	}
	acc.counter -= toSub
	return nil
}

func (acc *accumulator) numConvErr(num string) error {
	return fmt.Errorf("Couldn't convert %q to number", num)
}

func (acc *accumulator) loadLines() error {
	file, err := os.Open(acc.filename)
	if err != nil {
		return fmt.Errorf("Could not open file %s: %v", acc.filename, err)
	}
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			return fmt.Errorf("Ecountered blank line at line %d", lineNum)
		}
		acc.lines = append(acc.lines, line)
		lineNum++
	}

	return nil
}
