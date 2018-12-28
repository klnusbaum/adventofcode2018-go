package main

// start Dec 27 15:20
// end Dec 27 16:13

import (
	"fmt"
	"github.com/klnusbaum/adventofcode2018/ll"
	"os"
	"strconv"
	"strings"
)

var _allInsts map[string]instFunc = map[string]instFunc{
	"addr": addr,
	"addi": addi,
	"mulr": mulr,
	"muli": muli,
	"banr": banr,
	"bani": bani,
	"borr": borr,
	"bori": bori,
	"setr": setr,
	"seti": seti,
	"gtir": gtir,
	"gtri": gtri,
	"gtrr": gtrr,
	"eqir": eqir,
	"eqri": eqri,
	"eqrr": eqrr,
}

type instFunc func(regs regValues, inst instruction)

type regValues []int

type instruction struct {
	opCode  int
	inputA  int
	inputB  int
	outputC int
}

type sample struct {
	before regValues
	inst   instruction
	after  regValues
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

	samples := loadSamples(lines)

	threeOrMore := 0
	for _, s := range samples {
		if satisfiesThreeOrMore(s) {
			threeOrMore++
		}
	}

	fmt.Printf("Number of samples that fit three or more: %d\n", threeOrMore)

	return nil
}

func satisfiesThreeOrMore(s sample) bool {
	numSatisfy := 0

	for _, instF := range _allInsts {
		tempReg := make(regValues, 4)
		copy(tempReg, s.before)
		instF(tempReg, s.inst)
		if equalReg(tempReg, s.after) {
			numSatisfy++
		}

		if numSatisfy == 3 {
			return true
		}

	}
	return false
}

func equalReg(regs1, regs2 regValues) bool {
	for i, val := range regs1 {
		if regs2[i] != val {
			return false
		}
	}
	return true
}

func loadSamples(lines []string) []sample {
	curLine := 0
	var samples []sample
	for lines[curLine] != "" {
		newSample := sample{
			before: loadBeforeRegValues(lines[curLine]),
			inst:   loadInstruction(lines[curLine+1]),
			after:  loadAfterRegValues(lines[curLine+2]),
		}
		samples = append(samples, newSample)
		curLine += 4
	}

	return samples
}

func loadBeforeRegValues(line string) regValues {
	removedBefore := strings.TrimPrefix(line, "Before: ")
	return loadRegValues(removedBefore)
}

func loadAfterRegValues(line string) regValues {
	removedAfter := strings.TrimPrefix(line, "After: ")
	return loadRegValues(removedAfter)
}

func loadRegValues(line string) regValues {
	scrubbed := strings.TrimSpace(line)
	array := strings.TrimSuffix(strings.TrimPrefix(scrubbed, "["), "]")
	noCommas := strings.Replace(array, ",", "", -1)
	nums := strings.Split(noCommas, " ")
	reg0, _ := strconv.Atoi(nums[0])
	reg1, _ := strconv.Atoi(nums[1])
	reg2, _ := strconv.Atoi(nums[2])
	reg3, _ := strconv.Atoi(nums[3])
	return regValues{
		reg0,
		reg1,
		reg2,
		reg3,
	}
}

func loadInstruction(line string) instruction {
	nums := strings.Split(line, " ")
	opCode, _ := strconv.Atoi(nums[0])
	inputA, _ := strconv.Atoi(nums[1])
	inputB, _ := strconv.Atoi(nums[2])
	outputC, _ := strconv.Atoi(nums[3])

	return instruction{
		opCode:  opCode,
		inputA:  inputA,
		inputB:  inputB,
		outputC: outputC,
	}
}

func addr(regs regValues, inst instruction) {
	regs[inst.outputC] = regs[inst.inputA] + regs[inst.inputB]
}

func addi(regs regValues, inst instruction) {
	regs[inst.outputC] = regs[inst.inputA] + inst.inputB
}

func mulr(regs regValues, inst instruction) {
	regs[inst.outputC] = regs[inst.inputA] * regs[inst.inputB]
}

func muli(regs regValues, inst instruction) {
	regs[inst.outputC] = regs[inst.inputA] * inst.inputB
}

func banr(regs regValues, inst instruction) {
	regs[inst.outputC] = regs[inst.inputA] & regs[inst.inputB]
}

func bani(regs regValues, inst instruction) {
	regs[inst.outputC] = regs[inst.inputA] & inst.inputB
}

func borr(regs regValues, inst instruction) {
	regs[inst.outputC] = regs[inst.inputA] | regs[inst.inputB]
}

func bori(regs regValues, inst instruction) {
	regs[inst.outputC] = regs[inst.inputA] | inst.inputB
}

func setr(regs regValues, inst instruction) {
	regs[inst.outputC] = regs[inst.inputA]
}

func seti(regs regValues, inst instruction) {
	regs[inst.outputC] = inst.inputA
}

func gtir(regs regValues, inst instruction) {
	if inst.inputA > regs[inst.inputB] {
		regs[inst.outputC] = 1
	} else {
		regs[inst.outputC] = 0
	}
}

func gtri(regs regValues, inst instruction) {
	if regs[inst.inputA] > inst.inputB {
		regs[inst.outputC] = 1
	} else {
		regs[inst.outputC] = 0
	}
}

func gtrr(regs regValues, inst instruction) {
	if regs[inst.inputA] > regs[inst.inputB] {
		regs[inst.outputC] = 1
	} else {
		regs[inst.outputC] = 0
	}
}

func eqir(regs regValues, inst instruction) {
	if inst.inputA == regs[inst.inputB] {
		regs[inst.outputC] = 1
	} else {
		regs[inst.outputC] = 0
	}
}

func eqri(regs regValues, inst instruction) {
	if regs[inst.inputA] == inst.inputB {
		regs[inst.outputC] = 1
	} else {
		regs[inst.outputC] = 0
	}
}

func eqrr(regs regValues, inst instruction) {
	if regs[inst.inputA] == regs[inst.inputB] {
		regs[inst.outputC] = 1
	} else {
		regs[inst.outputC] = 0
	}
}
