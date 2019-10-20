package client

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Surge struct {
	UrlFilePath string
	Random      bool
}

func (surge Surge) Run() error {
	if surge.UrlFilePath != "" {
		file, err := os.Open(surge.UrlFilePath)
		if err != nil {
			return err
		}
		lines := []string{}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
		}

		if surge.Random {
			//https://yourbasic.org/golang/shuffle-slice-array/
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })
		}

		for _, line := range lines {
			args := strings.Fields(line)
			HttpCommand{}.Execute(args)
		}

		if err := scanner.Err(); err != nil {
			return err
		}
	}
	return nil
}
