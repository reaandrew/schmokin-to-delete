package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func CreateRandomLines(lineCount int) []string {
	lines := []string{}
	for i := 0; i < lineCount; i++ {
		lines = append(lines, fmt.Sprintf("http://localhost:8080/%v", i+1))
	}
	return lines
}

func CreateRandomHTTPTestFile(lineCount int) *os.File {
	lines := CreateRandomLines(lineCount)
	return CreateTestFile(lines)
}

func CreateTestFile(lines []string) *os.File {
	fileContents := strings.Join(lines, "\n")
	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(file.Name(), []byte(fileContents), os.ModePerm)
	if err != nil {
		panic(err)
	}
	return file
}

func MapStrings(array []string, delegate func(value string) string) (values []string) {
	for _, value := range array {
		values = append(values, delegate(value))
	}
	return
}
