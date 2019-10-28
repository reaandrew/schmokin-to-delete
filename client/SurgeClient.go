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
	UrlFilePath  string
	Random       bool
	WorkerCount  int
	Iterations   int
	HttpClient   HttpClient
	lock         sync.Mutex
	waitGroup    sync.WaitGroup
	transactions int
	errors       int
}

func (surge *Surge) execute(lines []string) Result {
	var lock = sync.Mutex{}
	transactions := 0
	errors := 0
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
				err := command.Execute(args)
				lock.Lock()
				if err != nil {
					errors++
				}
				transactions++
				lock.Unlock()
				if i > 0 && i == surge.Iterations-1 {
					break
				}
			}
			wg.Done()
		}(lines)
	}
	wg.Wait()
	result := Result{
		Transactions: transactions,
	}
	if errors == 0 {
		result.Availability = 1
	} else {
		result.Availability = float64(1 - float64(errors)/float64(transactions))
	}
	return result
}

func (surge Surge) Run() (result Result, err error) {
	var file *os.File
	if surge.UrlFilePath != "" {
		file, err = os.Open(surge.UrlFilePath)
		if err != nil {
			return
		}
		lines := []string{}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
		}
		if err = scanner.Err(); err != nil {
			return
		}

		if surge.Random {
			//https://yourbasic.org/golang/shuffle-slice-array/
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })
		}

		result = surge.execute(lines)
	}
	return
}
