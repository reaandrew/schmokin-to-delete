package utils

import (
	"bufio"
	"os"
)

func ReadFileToLines(path string) (lines []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	lines = []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	err = scanner.Err()
	return
}
