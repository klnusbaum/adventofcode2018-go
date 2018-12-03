package ll

import (
	"bufio"
	"fmt"
	"os"
)

type LineLoader struct {
	filename string
}

func NewLineLoader(filename string) LineLoader {
	return LineLoader{
		filename: filename,
	}
}

func (loader LineLoader) Load() ([]string, error) {
	file, err := os.Open(loader.filename)
	if err != nil {
		return nil, fmt.Errorf("Couldn't open file %q: %v", loader.filename, err)
	}

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}
