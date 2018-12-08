package main

/*
 */

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	if err := analyze("input.txt"); err != nil {
		fmt.Printf("Error analyzing: %v\n", err)
		os.Exit(1)
	}
}

func analyze(filename string) error {
	nums, err := getNums(filename)
	if err != nil {
		return fmt.Errorf("error geting nums: %v", err)
	}

	val := getVal(nums)

	fmt.Printf("Val %d\n", val)

	return nil
}

func getVal(q *queue) int {
	numChildren := q.deque()
	numMetadata := q.deque()
	childVals := make([]int, numChildren)
	for i := 0; i < numChildren; i++ {
		childVals[i] = getVal(q)
	}

	metadataVals := make([]int, numMetadata)
	for i := 0; i < numMetadata; i++ {
		metadataVals[i] = q.deque()
	}

	if numChildren == 0 {
		sum := 0
		for _, val := range metadataVals {
			sum += val
		}
		return sum
	}

	val := 0
	for _, childIndex := range metadataVals {
		actualIndex := childIndex - 1
		if actualIndex >= 0 && actualIndex < len(childVals) {
			val += childVals[actualIndex]
		}
	}

	return val
}

func getNums(filename string) (*queue, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Couldn't open file %s: %v", filename, err)
	}

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("Error loading file: %v", err)
	}
	strNums := strings.Split(strings.TrimSpace(string(contents)), " ")

	nums := make([]int, len(strNums))

	for i, entry := range strNums {
		num, err := strconv.Atoi(entry)
		if err != nil {
			return nil, fmt.Errorf("Error converting number %s: %v", entry, err)
		}
		nums[i] = num
	}

	return &queue{
		ints: nums,
	}, nil

}

type queue struct {
	ints []int
}

func (q *queue) deque() int {
	toReturn := q.ints[0]
	q.ints = q.ints[1:]
	return toReturn
}

func (q *queue) peek() int {
	return q.ints[0]
}

func (q *queue) isEmpty() bool {
	return len(q.ints) == 0
}
