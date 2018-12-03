package main

import (
	"testing"
)

func TestDistance(t *testing.T) {
	tests := []struct {
		msg  string
		str1 string
		str2 string
		dist int
	}{
		{
			msg:  "no distance",
			str1: "hello",
			str2: "hello",
			dist: 0,
		},
		{
			msg:  "distance 1",
			str1: "hello",
			str2: "hellg",
			dist: 1,
		},
		{
			msg:  "distance 2",
			str1: "hello",
			str2: "hillg",
			dist: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			dist := distance(tt.str1, tt.str2)
			if dist != tt.dist {
				t.Errorf("Expected distance between %q and %q was %d, instead got %d", tt.str1, tt.str2, tt.dist, dist)
			}
		})
	}

}
