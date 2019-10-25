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
	Iterations  int
	HttpClient  HttpClient
}

func (surge Surge) execute(lines []string) int {
	transactions := 0
	var wg sync.WaitGroup
	for i := 0; i < surge.WorkerCount; i++ {
		wg.Add(1)
		go func(linesValue []string) {
			for i := 0; i < len(linesValue) || (surge.Iterations > 0 && i < surge.Iterations); i++ {
				line := linesValue[i%len(linesValue)]
				var command = HttpCommand{
					client: surge.HttpClient,
				}
				var args = strings.Fields(line)
				command.Execute(args)
			}
			wg.Done()
		}(lines)
	}
	wg.Wait()
	return transactions
}

func (surge Surge) Run() (int, error) {
	transactions := 0
	if surge.UrlFilePath != "" {
		file, err := os.Open(surge.UrlFilePath)
		if err != nil {
			return transactions, err
		}
		lines := []string{}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
		}
		if err := scanner.Err(); err != nil {
			return transactions, err
		}

		if surge.Random {
			//https://yourbasic.org/golang/shuffle-slice-array/
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })
		}

		transactions = surge.execute(lines)
	}
	return transactions, nil
}
