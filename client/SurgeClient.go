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

func (surge *Surge) worker(linesValue []string) {
	for i := 0; i < len(linesValue) || (surge.Iterations > 0 && i < surge.Iterations); i++ {
		line := linesValue[i%len(linesValue)]
		var command = HttpCommand{
			client: surge.HttpClient,
		}
		var args = strings.Fields(line)
		err := command.Execute(args)
		surge.lock.Lock()
		if err != nil {
			surge.errors++
		}
		surge.transactions++
		surge.lock.Unlock()
		if i > 0 && i == surge.Iterations-1 {
			break
		}
	}
	surge.waitGroup.Done()
}

func (surge *Surge) execute(lines []string) Result {
	for i := 0; i < surge.WorkerCount; i++ {
		surge.waitGroup.Add(1)
		go surge.worker(lines)
	}
	surge.waitGroup.Wait()
	result := Result{
		Transactions: surge.transactions,
	}
	if surge.errors == 0 {
		result.Availability = 1
	} else {
		result.Availability = float64(1 - float64(surge.errors)/float64(surge.transactions))
	}
	return result
}

func (surge *Surge) Run() (result Result, err error) {
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
