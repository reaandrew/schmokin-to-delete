package client

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type Surge struct {
	UrlFilePath string
	Random      bool
	WorkerCount int
}

func (surge Surge) execute(lines []string) {
	var wg sync.WaitGroup
	for i := 0; i < surge.WorkerCount; i++ {
		wg.Add(1)
		go func(linesValue []string) {
			for _, line := range linesValue {
				var command = HttpCommand{}
				var args = strings.Fields(line)
				command.Execute(args)
			}
			wg.Done()
		}(lines)
	}
	wg.Wait()
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

		surge.execute(lines)

		if err := scanner.Err(); err != nil {
			return err
		}
	}
	return nil
}
