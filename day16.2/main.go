package main

// start Dec 27 15:20
// end Dec 27 16:13
// end Dec 27 17:16

import (
	"fmt"
	"github.com/klnusbaum/adventofcode2018/ll"
	"os"
	"strconv"
	"strings"
)

const (
	_progStart = 3166
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
	actualOpCodes := findOpCodes(samples)
	program := loadProgram(lines)

	realRegs := regValues{0, 0, 0, 0}
	runProgram(realRegs, program, actualOpCodes)

	fmt.Printf("Value of register 0: %d\n", realRegs[0])

	return nil
}

func runProgram(regs regValues, program []instruction, opCodes map[int]instFunc) {
	for _, inst := range program {
		instF := opCodes[inst.opCode]
		instF(regs, inst)
	}
}

func findOpCodes(samples []sample) map[int]instFunc {
	possibilities := initialPossibilities()
	filterPossibiliesWithSamples(possibilities, samples)
	constrainPossibilities(possibilities)

	realOpcodes := make(map[int]instFunc)
	for opCode, poss := range possibilities {
		if len(poss) != 1 {
			panic(fmt.Sprintf("Couldn't find inst for opcode %d, lenth of poss was %d", opCode, len(poss)))
		}

		for _, instF := range poss {
			realOpcodes[opCode] = instF
		}
	}

	return realOpcodes
}

func constrainPossibilities(possibilities map[int]map[string]instFunc) {
	for notFullyConstrained(possibilities) {
		for opCode, poss := range possibilities {
			if len(poss) != 1 {
				continue
			}

			constainedOpName := ""
			for opName, _ := range poss {
				constainedOpName = opName
			}

			removeOpNameFromOpCodes(constainedOpName, opCode, possibilities)
		}
	}
}

func removeOpNameFromOpCodes(opName string, realOpCode int, possibilities map[int]map[string]instFunc) {
	for opCode, poss := range possibilities {
		if opCode == realOpCode {
			continue
		}
		delete(poss, opName)
	}
}

func notFullyConstrained(possibilities map[int]map[string]instFunc) bool {
	for _, poss := range possibilities {
		if len(poss) > 1 {
			return true
		}
	}

	return false
}

func filterPossibiliesWithSamples(possibilities map[int]map[string]instFunc, samples []sample) {
	for _, s := range samples {
		for opName, instF := range _allInsts {
			tempReg := make(regValues, 4)
			copy(tempReg, s.before)
			instF(tempReg, s.inst)
			if !equalReg(tempReg, s.after) {
				delete(possibilities[s.inst.opCode], opName)
			}
		}
	}
}

func initialPossibilities() map[int]map[string]instFunc {
	possibilities := make(map[int]map[string]instFunc)
	for i := 0; i < 16; i++ {
		poss := make(map[string]instFunc)
		for name, instF := range _allInsts {
			poss[name] = instF
		}
		possibilities[i] = poss
	}

	return possibilities
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

func loadProgram(lines []string) []instruction {
	curLine := _progStart
	var insts []instruction
	for curLine < len(lines) {
		insts = append(insts, loadInstruction(lines[curLine]))
		curLine++
	}
	return insts
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
