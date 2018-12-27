package main

// start Dec 16 6:17 am
// finish
/*
 */

import (
	"fmt"
	"github.com/klnusbaum/adventofcode2018/ll"
	"os"
	"sort"
)

const (
	_effectiveInf = 1000000000
	_infHP        = 10000
)

type game struct {
	tiles     [][]tileType
	openTiles map[point]bool
	creatures map[point]creature
}

func main() {
	g1 := &game{}
	if err := g1.analyze("simple4.txt"); err != nil {
		fmt.Printf("Error analyzing: %v\n", err)
		os.Exit(1)
	}
}

func (g *game) analyze(filename string) error {
	loader := ll.NewLineLoader(filename)
	lines, err := loader.Load()
	if err != nil {
		return fmt.Errorf("Couldn't load lines: %v\n", err)
	}

	g.loadMapAndCreatures(lines)

	fmt.Println("Initial")
	g.printMap()

	iterations := 0
	for {
		creaturesToIterate := g.sortedCreatures()
		middleOfRoundBreak := false
		for _, p := range creaturesToIterate {
			creature, ok := g.creatures[p]
			if !ok {
				continue
			}

			targets := g.getTargets(creature.ctype)
			if len(targets) == 0 {
				middleOfRoundBreak = true
				break
			}

			didAttack := g.maybeAttack(creature, p)
			if didAttack {
				continue
			}

			newPos := g.maybeMove(targets, p)
			if newPos.x != -1 && newPos.y != -1 {
				g.maybeAttack(creature, newPos)
			}
		}

		if !middleOfRoundBreak {
			iterations++
			fmt.Printf("After %d\n", iterations)
			g.printMap()
		}

		if len(g.getTargets(elf)) == 0 || len(g.getTargets(goblin)) == 0 {
			break
		}
	}

	hpSum := g.getHPSum()
	fmt.Printf("Outcome %d * %d: %d\n", iterations, hpSum, hpSum*iterations)

	return nil
}

func (g *game) sortedCreatures() points {
	var cpoints points
	for p, _ := range g.creatures {
		cpoints = append(cpoints, p)
	}

	sort.Sort(cpoints)
	return cpoints
}

func (g *game) getTargets(forType ctype) map[point]creature {
	targets := make(map[point]creature)
	targetType := elf
	if forType == elf {
		targetType = goblin
	}
	for p, creature := range g.creatures {
		if creature.ctype == targetType {
			targets[p] = creature
		}
	}
	return targets
}

func (g *game) maybeAttack(attacker creature, attackerPos point) bool {
	ncPositions := g.neighboorCreatures(attacker, attackerPos)
	if len(ncPositions) == 0 {
		return false
	}

	smallestHP := _infHP
	var smallestToAttack points
	for _, creaturePos := range ncPositions {
		creatureHP := g.creatures[creaturePos].hp
		if creatureHP < smallestHP {
			smallestHP = creatureHP
			smallestToAttack = points{creaturePos}
		} else if creatureHP == smallestHP {
			smallestToAttack = append(smallestToAttack, creaturePos)
		}
	}

	sort.Sort(smallestToAttack)
	g.attack(smallestToAttack[0])
	return true
}

func (g *game) attack(pos point) {
	creatureToAttack, ok := g.creatures[pos]
	if !ok {
		panic("Attacking creature that doesn't exist")
	}

	newCreature := creature{
		ctype: creatureToAttack.ctype,
		hp:    creatureToAttack.hp - 3,
	}

	if newCreature.hp <= 0 {
		delete(g.creatures, pos)
	} else {
		g.creatures[pos] = newCreature
	}
}

func (g *game) maybeMove(targets map[point]creature, moverPos point) point {
	closestOpen := g.findClosestOpen(targets, moverPos)
	if len(closestOpen) == 0 {
		return point{-1, -1}
	}

	sort.Sort(closestOpen)
	selectedOpen := closestOpen[0]
	return g.moveTowardsOpen(moverPos, selectedOpen)
}

func (g *game) findClosestOpen(targets map[point]creature, moverPos point) points {
	var closestInRange points
	shortestDistance := _effectiveInf

	for p, _ := range g.inRange(targets) {
		distToOpen := g.getDistTo(moverPos, p)
		if distToOpen == -1 {
			continue
		}
		if distToOpen < shortestDistance {
			shortestDistance = distToOpen
			closestInRange = points{p}
		} else if distToOpen == shortestDistance {
			closestInRange = append(closestInRange, p)
		}
	}

	return closestInRange
}

func (g *game) inRange(targets map[point]creature) map[point]bool {
	inRange := make(map[point]bool)
	for p, _ := range targets {
		for _, neighboor := range g.getOpenNeighboors(p) {
			inRange[neighboor] = true
		}
	}

	return inRange
}

func (g *game) getDistTo(moverPos, targetPos point) int {
	unvisited := make(map[point]bool)
	distanceTo := make(map[point]int)
	for p, _ := range g.openTiles {
		if _, ok := g.creatures[p]; !ok {
			unvisited[p] = true
			distanceTo[p] = _effectiveInf
		}
	}

	unvisited[moverPos] = true
	distanceTo[moverPos] = 0
	currentNode := moverPos
	for unvisited[targetPos] && reachableUnivisited(unvisited, distanceTo) {
		for _, neighboor := range g.unvisitedNeighboors(unvisited, currentNode) {
			tenativeDistance := distanceTo[currentNode] + 1
			if tenativeDistance < distanceTo[neighboor] {
				distanceTo[neighboor] = tenativeDistance
			}
		}

		delete(unvisited, currentNode)
		currentNode = smallestUnvisited(unvisited, distanceTo)
	}

	if _, ok := unvisited[targetPos]; !ok {
		return distanceTo[targetPos]
	}

	return -1
}

func reachableUnivisited(unvisited map[point]bool, distanceTo map[point]int) bool {
	for p, _ := range unvisited {
		if distanceTo[p] != _effectiveInf {
			return true
		}
	}

	return false
}

func (g *game) unvisitedNeighboors(unvisited map[point]bool, pos point) points {
	var unvisitedNeighboors points
	for _, p := range g.getOpenNeighboors(pos) {
		if _, ok := unvisited[p]; ok {
			unvisitedNeighboors = append(unvisitedNeighboors, p)
		}
	}

	return unvisitedNeighboors
}

func smallestUnvisited(unvisited map[point]bool, distanceTo map[point]int) point {
	smallest := _effectiveInf
	p := point{}

	for p1, _ := range unvisited {
		if distanceTo[p1] < smallest {
			smallest = distanceTo[p1]
			p = p1
		}
	}

	return p
}

func (g *game) moveTowardsOpen(moverPos point, target point) point {
	shortestDist := _effectiveInf
	var shortestNeighboors points

	for _, neighboor := range g.getOpenNeighboors(moverPos) {
		dist := g.getDistTo(neighboor, target)
		if dist == -1 {
			continue
		} else if dist < shortestDist {
			shortestDist = dist
			shortestNeighboors = points{neighboor}
		} else if dist == shortestDist {
			shortestNeighboors = append(shortestNeighboors, neighboor)
		}
	}

	if shortestDist == _effectiveInf {
		panic("moving a point toward a target, but no path to target")
	}

	sort.Sort(shortestNeighboors)

	moveTo := shortestNeighboors[0]
	movingCreature := g.creatures[moverPos]

	delete(g.creatures, moverPos)
	g.creatures[moveTo] = movingCreature
	return moveTo
}

func (g *game) getOpenNeighboors(p point) points {
	above := point{p.x, p.y - 1}
	below := point{p.x, p.y + 1}
	left := point{p.x - 1, p.y}
	right := point{p.x + 1, p.y}

	potentials := points{above, below, left, right}
	var realNeighboords points
	for _, potential := range potentials {
		_, open := g.openTiles[potential]
		_, creature := g.creatures[potential]
		if open && !creature {
			realNeighboords = append(realNeighboords, potential)
		}
	}

	return realNeighboords
}

func (g *game) neighboorCreatures(cr creature, p point) points {
	above := point{p.x, p.y - 1}
	below := point{p.x, p.y + 1}
	left := point{p.x - 1, p.y}
	right := point{p.x + 1, p.y}

	potentials := points{above, below, left, right}
	var realNeighboors points
	for _, potential := range potentials {
		creature, okCreature := g.creatures[potential]
		if okCreature && cr.ctype != creature.ctype {
			realNeighboors = append(realNeighboors, potential)
		}
	}
	return realNeighboors
}

func (g *game) getHPSum() int {
	sum := 0
	for _, creature := range g.creatures {
		sum += creature.hp
	}
	return sum
}

func (g *game) loadMapAndCreatures(lines []string) {
	g.tiles = make([][]tileType, len(lines))
	g.creatures = make(map[point]creature)
	g.openTiles = make(map[point]bool)

	for j, line := range lines {
		g.tiles[j] = make([]tileType, len(line))
		for i, char := range line {
			p := point{i, j}
			if char == '#' {
				g.tiles[j][i] = wall
			} else if char == '.' {
				g.tiles[j][i] = open
				g.openTiles[p] = true
			} else if char == 'G' {
				g.tiles[j][i] = open
				g.creatures[p] = creature{goblin, 200}
				g.openTiles[p] = true
			} else if char == 'E' {
				g.tiles[j][i] = open
				g.creatures[p] = creature{elf, 200}
				g.openTiles[p] = true
			} else {
				panic(fmt.Sprintf("Unrecognized input char: %c", char))
			}
		}
	}
}

func (g *game) printMap() {
	for j, rowTiles := range g.tiles {
		var cinfo []creature
		for i, tile := range rowTiles {
			if creature, ok := g.creatures[point{i, j}]; ok {
				cinfo = append(cinfo, creature)
				if creature.ctype == goblin {
					fmt.Print("G")
				} else if creature.ctype == elf {
					fmt.Print("E")
				} else {
					panic(fmt.Sprintf("Unknown creature type: %d", creature.ctype))
				}
				continue
			}

			if tile == open {
				fmt.Print(".")
			} else if tile == wall {
				fmt.Print("#")
			} else {
				panic(fmt.Sprintf("Unknown tile type: %d", tile))
			}
		}
		fmt.Print("  ")
		for _, creature := range cinfo {
			fmt.Print(" ")
			if creature.ctype == goblin {
				fmt.Print("G")
			} else {
				fmt.Print("E")
			}
			fmt.Printf("(%d)", creature.hp)
		}
		fmt.Println()
	}
}

type tileType int

const (
	open tileType = iota
	wall
)

type ctype int

const (
	elf ctype = iota
	goblin
)

type creature struct {
	ctype ctype
	hp    int
}

type point struct {
	x int
	y int
}

type points []point

func (p points) Len() int      { return len(p) }
func (p points) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p points) Less(i, j int) bool {
	if p[i].y < p[j].y {
		return true
	} else if p[i].y == p[j].y && p[i].x < p[j].x {
		return true
	}

	return false
}
