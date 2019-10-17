package client

import (
	"bufio"
	"os"
	"strings"
)

type Surge struct {
	UrlFilePath string
}

func (surge Surge) Run() error {
	if surge.UrlFilePath != "" {
		file, err := os.Open(surge.UrlFilePath)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			args := strings.Fields(line)
			HttpCommand{}.Execute(args)
		}

		if err := scanner.Err(); err != nil {
			return err
		}
	}
	return nil
}
