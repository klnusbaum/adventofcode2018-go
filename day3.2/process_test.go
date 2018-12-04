package main

import (
	"testing"
)

func TestLoad(t *testing.T) {
	claims, _ := loadClaims("input.txt")

	if claims[0].x != 604 {
		t.Errorf("x doesn't match")
	}

	if claims[0].y != 100 {
		t.Errorf("y doesn't match")
	}

	if claims[0].w != 17 {
		t.Errorf("w doesn't match")
	}

	if claims[0].h != 27 {
		t.Errorf("h doesn't match")
	}

	if claims[0].id != 1 {
		t.Errorf("id doesn't match")
	}
}

func TestProcess(t *testing.T) {
	claim1 := claim{
		x:  1,
		y:  3,
		w:  4,
		h:  4,
		id: 1,
	}

	claim2 := claim{
		x:  3,
		y:  1,
		w:  4,
		h:  4,
		id: 2,
	}

	claim3 := claim{
		x:  5,
		y:  5,
		w:  2,
		h:  2,
		id: 3,
	}

	claims := []claim{claim1, claim2, claim3}
	sheet, _ := processClaims(claims)

	sequals(t, sheet[0][0], []int{})
	sequals(t, sheet[1][3], []int{2})
	sequals(t, sheet[3][3], []int{1, 2})
	sequals(t, sheet[4][3], []int{1, 2})
	sequals(t, sheet[5][5], []int{3})
}

func TestFindSingle(t *testing.T) {
	claim1 := claim{
		x:  1,
		y:  3,
		w:  4,
		h:  4,
		id: 1,
	}

	claim2 := claim{
		x:  3,
		y:  1,
		w:  4,
		h:  4,
		id: 2,
	}

	claim3 := claim{
		x:  5,
		y:  5,
		w:  2,
		h:  2,
		id: 3,
	}

	claims := []claim{claim1, claim2, claim3}
	sheet, _ := processClaims(claims)

	id := findSingleId(sheet)
	if id != 3 {
		t.Errorf("Id was not 3, it was %d", id)
	}
}

func sequals(t *testing.T, s1, s2 []int) {
	if len(s1) != len(s2) {
		t.Errorf("Slices with different length: %v, %v", s1, s2)
	}

	for i, val := range s1 {
		if val != s2[i] {
			t.Errorf("Mismatch at position %d: %v and %v", i, s1, s2)
		}
	}
}
