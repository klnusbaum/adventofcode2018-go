package main

/*

--- Day 4: Repose Record ---
You've sneaked into another supply closet - this time, it's across from the prototype suit manufacturing lab. You need to sneak inside and fix the issues with the suit, but there's a guard stationed outside the lab, so this is as close as you can safely get.

As you search the closet for anything that might help, you discover that you're not the first person to want to sneak in. Covering the walls, someone has spent an hour starting every midnight for the past few months secretly observing this guard post! They've been writing down the ID of the one guard on duty that night - the Elves seem to have decided that one guard was enough for the overnight shift - as well as when they fall asleep or wake up while at their post (your puzzle input).

For example, consider the following records, which have already been organized into chronological order:

[1518-11-01 00:00] Guard #10 begins shift
[1518-11-01 00:05] falls asleep
[1518-11-01 00:25] wakes up
[1518-11-01 00:30] falls asleep
[1518-11-01 00:55] wakes up
[1518-11-01 23:58] Guard #99 begins shift
[1518-11-02 00:40] falls asleep
[1518-11-02 00:50] wakes up
[1518-11-03 00:05] Guard #10 begins shift
[1518-11-03 00:24] falls asleep
[1518-11-03 00:29] wakes up
[1518-11-04 00:02] Guard #99 begins shift
[1518-11-04 00:36] falls asleep
[1518-11-04 00:46] wakes up
[1518-11-05 00:03] Guard #99 begins shift
[1518-11-05 00:45] falls asleep
[1518-11-05 00:55] wakes up
Timestamps are written using year-month-day hour:minute format. The guard falling asleep or waking up is always the one whose shift most recently started. Because all asleep/awake times are during the midnight hour (00:00 - 00:59), only the minute portion (00 - 59) is relevant for those events.

Visually, these records show that the guards are asleep at these times:

Date   ID   Minute
            000000000011111111112222222222333333333344444444445555555555
            012345678901234567890123456789012345678901234567890123456789
11-01  #10  .....####################.....#########################.....
11-02  #99  ........................................##########..........
11-03  #10  ........................#####...............................
11-04  #99  ....................................##########..............
11-05  #99  .............................................##########.....
The columns are Date, which shows the month-day portion of the relevant day; ID, which shows the guard on duty that day; and Minute, which shows the minutes during which the guard was asleep within the midnight hour. (The Minute column's header shows the minute's ten's digit in the first row and the one's digit in the second row.) Awake is shown as ., and asleep is shown as #.

Note that guards count as asleep on the minute they fall asleep, and they count as awake on the minute they wake up. For example, because Guard #10 wakes up at 00:25 on 1518-11-01, minute 25 is marked as awake.

If you can figure out the guard most likely to be asleep at a specific time, you might be able to trick that guard into working tonight so you can have the best chance of sneaking in. You have two strategies for choosing the best guard/minute combination.

Strategy 1: Find the guard that has the most minutes asleep. What minute does that guard spend asleep the most?

In the example above, Guard #10 spent the most minutes asleep, a total of 50 minutes (20+25+5), while Guard #99 only slept for a total of 30 minutes (10+10+10). Guard #10 was asleep most during minute 24 (on two days, whereas any other minute the guard was asleep was only seen on one day).

While this example listed the entries in chronological order, your entries are in the order you found them. You'll need to organize them before they can be analyzed.

What is the ID of the guard you chose multiplied by the minute you chose? (In the above example, the answer would be 10 * 24 = 240.)
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

	guard := mostAsleep(gs)
	fmt.Printf("Most asleep guard: %s\n", guard)

	min := mostAsleepMinute(gs[guard])

	fmt.Printf("Minute most likely to be asleep: %d\n", min)

	guardInt, _ := strconv.Atoi(strings.TrimLeft(guard, "#"))

	fmt.Printf("Most asleep minute * guard id: %d\n", min*guardInt)

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

func mostAsleep(gs guardSchedules) string {
	sleepiestGuard := "#notfound"
	maxAsleep := 0
	for guard, schedules := range gs {
		asleep := totalTimeAsleep(schedules)
		if asleep > maxAsleep {
			maxAsleep = asleep
			sleepiestGuard = guard
		}
	}
	return sleepiestGuard
}

func totalTimeAsleep(schedules schedules) int {
	totalTime := 0
	for _, schedule := range schedules {
		totalTime += timeForSchedule(schedule)
	}
	return totalTime
}

func timeForSchedule(actions actions) int {
	minsAsleep := 0
	state := "up"
	for i := 0; i < 60; i++ {
		if action, ok := actions[i]; ok {
			state = action
		}
		if state == "asleep" {
			minsAsleep++
		}
	}
	return minsAsleep
}

func mostAsleepMinute(schedules schedules) int {
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

func histoMax(histo map[int]int) int {
	max := 0
	maxMinute := 0
	for minute, num := range histo {
		if num > max {
			max = num
			maxMinute = minute
		}
	}

	return maxMinute
}
