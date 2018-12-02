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
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Couldn't open file %q: %v", filename, err)
		os.Exit(1)
	}

	acc := newAccumulator(file)

	if err := acc.accumulate(os.Stdout); err != nil {
		fmt.Printf("Couldn't accumulate result: %v", err)
		os.Exit(1)
	}
}

func newAccumulator(r io.Reader) *accumulator {
	return &accumulator{
		scanner: bufio.NewScanner(r),
		counter: 0,
		lineNum: 0,
	}
}

type accumulator struct {
	scanner *bufio.Scanner
	counter int
	lineNum int
}

func (acc *accumulator) accumulate(out io.Writer) error {
	for acc.scanner.Scan() {
		line := acc.scanner.Text()
		if err := acc.accLine(line); err != nil {
			return fmt.Errorf("Couldn't accumulate line %q: %v", err)
		}
		acc.lineNum++
	}

	_, err := fmt.Fprintf(out, "Accumulated Value: %d\n", acc.counter)
	return err
}

func (acc *accumulator) accLine(line string) error {
	if len(line) <= 0 {
		return fmt.Errorf("Line %d is blank", acc.lineNum)
	}
	switch firstChar := string(line[0]); firstChar {
	case "+":
		return acc.addNum(line[1:])
	case "-":
		return acc.subNum(line[1:])
	default:
		return fmt.Errorf("Unrecognized first character %q on line %d", firstChar, acc.lineNum)
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
	return fmt.Errorf("Couldn't convert %q to number on line %d", num, acc.lineNum)
}
