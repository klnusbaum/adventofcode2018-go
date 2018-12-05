package main

/*
Strategy 2: Of all guards, which guard is most frequently asleep on the same minute?

In the example above, Guard #99 spent minute 45 asleep more than any other guard or minute - three times in total. (In all other cases, any guard spent any minute asleep at most twice.)

What is the ID of the guard you chose multiplied by the minute you chose? (In the above example, the answer would be 99 * 45 = 4455.)
*/

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/klnusbaum/adventofcode2018/ll"
)

// minute -> action (up or asleep)
type actions map[int]string

//day to actions on day
type schedules map[string]actions

// guard id to schedule
type guardSchedules map[string]schedules

func (gs guardSchedules) String() string {
	var sb strings.Builder
	for guard, schedules := range gs {

		sb.WriteString(fmt.Sprintf("Guard %s\n", guard))
		for day, actions := range schedules {
			sb.WriteString(fmt.Sprintf("Day %s: %v\n", day, actions))
		}
	}

	return sb.String()
}

func main() {
	filename := os.Args[1]

	if err := analyze(filename); err != nil {
		fmt.Printf("Error analyzing schedules: %v\n", err)
		os.Exit(1)
	}
}

func analyze(filename string) error {
	loader := ll.NewLineLoader(filename)
	lines, err := loader.Load()
	if err != nil {
		return fmt.Errorf("couldn't load lines: %v", err)
	}

	gs, err := loadGuards(lines)
	if err != nil {
		return fmt.Errorf("couldn't load guards: %v", err)
	}

	maxGuard, maxMinute := consistentlyAsleepGuard(gs)

	guardInt, _ := strconv.Atoi(strings.TrimLeft(maxGuard, "#"))

	fmt.Printf("Guard %s is asleep at minute %d very often\n", maxGuard, maxMinute)
	fmt.Printf("Those two number multipled are %d\n", guardInt*maxMinute)

	return nil
}

func loadGuards(lines []string) (guardSchedules, error) {
	gs := make(guardSchedules)
	sort.Strings(lines)
	currentGuard := ""
	for _, line := range lines {
		fields := strings.Fields(line)
		if strings.Contains(line, "Guard") {
			currentGuard = fields[3]
			continue
		}

		day := strings.TrimLeft(fields[0], "[")
		minute := strings.TrimRight(strings.Split(fields[1], ":")[1], "]")
		min, _ := strconv.Atoi(minute)
		aType := fields[3]

		guard := gs[currentGuard]
		if guard == nil {
			guard = make(schedules)
			gs[currentGuard] = guard
		}

		as := guard[day]
		if as == nil {
			as = make(map[int]string)
			guard[day] = as
		}

		as[min] = aType
	}

	return gs, nil
}

func consistentlyAsleepGuard(gs guardSchedules) (string, int) {
	maxGuard := "guard not found"
	maxMinute := 0
	maxMinuteCount := 0
	for guard, schedules := range gs {
		minute, count := mostAsleepMinute(schedules)
		fmt.Printf("Guard %s is asleep at %d, %d times\n", guard, minute, count)
		if count > maxMinuteCount {
			maxGuard = guard
			maxMinute = minute
			maxMinuteCount = count
		}
	}

	return maxGuard, maxMinute
}

func mostAsleepMinute(schedules schedules) (int, int) {
	minHistogram := make(map[int]int)
	for _, schedule := range schedules {
		state := "up"
		for i := 0; i < 60; i++ {
			if action, ok := schedule[i]; ok {
				state = action
			}
			if state == "asleep" {
				minHistogram[i] = minHistogram[i] + 1
			}
		}
	}

	return histoMax(minHistogram)
}

func histoMax(histo map[int]int) (int, int) {
	max := 0
	maxMinute := 0
	for minute, num := range histo {
		if num > max {
			max = num
			maxMinute = minute
		}
	}

	return maxMinute, max
}
