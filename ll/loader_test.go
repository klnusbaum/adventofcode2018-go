package ll

import (
	"testing"
)

const (
	_testFile = "testdata/simplefile1.txt"
)

var (
	_expectedLines = []string{"one", "two", "three"}
)

func TestLoadFile(t *testing.T) {
	loader := NewLineLoader("testdata/simplefile1.txt")
	lines, err := loader.Load()

	if err != nil {
		t.Errorf("Did not expect load error: %v", err)
	}

	for i := range lines {
		if lines[i] != _expectedLines[i] {
			t.Errorf("Non-matching line %d. Expected %q but got %q", i, _expectedLines[i], lines[i])
		}
	}
}
