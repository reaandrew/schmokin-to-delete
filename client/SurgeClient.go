package client

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/reaandrew/surge/utils"
)

type surge struct {
	//TODO: Create a configuration struct for these
	urlFilePath string
	random      bool
	workerCount int
	iterations  int
	httpClient  HttpClient
	timer       utils.Timer
	lock        sync.Mutex
	waitGroup   sync.WaitGroup
	//TODO: Create a stats struct for these
	transactions       int
	errors             int
	totalBytesSent     int
	totalBytesReceived int
	responseTime       metrics.Histogram
}

func (surge *surge) worker(linesValue []string) {
	for i := 0; i < len(linesValue) || (surge.iterations > 0 && i < surge.iterations); i++ {
		line := linesValue[i%len(linesValue)]
		var command = HttpCommand{
			client: surge.httpClient,
			timer:  surge.timer,
		}
		var args = strings.Fields(line)
		result := command.Execute(args)
		surge.lock.Lock()
		if result.Error != nil {
			surge.errors++
		}
		surge.transactions++
		surge.totalBytesSent += result.TotalBytesSent
		surge.totalBytesReceived += result.TotalBytesReceived
		surge.responseTime.Update(int64(result.ResponseTime))
		surge.lock.Unlock()
		if i > 0 && i == surge.iterations-1 {
			break
		}
	}
	surge.waitGroup.Done()
}

func (surge *surge) execute(lines []string) Result {
	for i := 0; i < surge.workerCount; i++ {
		surge.timer.Start()
		surge.waitGroup.Add(1)
		go surge.worker(lines)
	}
	surge.waitGroup.Wait()
	result := Result{
		Transactions:        surge.transactions,
		ElapsedTime:         surge.timer.Stop(),
		TotalBytesSent:      surge.totalBytesSent,
		TotalBytesReceived:  surge.totalBytesReceived,
		AverageResponseTime: surge.responseTime.Mean(),
	}
	if surge.errors == 0 {
		result.Availability = 1
	} else {
		result.Availability = float64(1 - float64(surge.errors)/float64(surge.transactions))
	}
	return result
}

func (surge *surge) Run() (result Result, err error) {
	var file *os.File
	if surge.urlFilePath != "" {
		file, err = os.Open(surge.urlFilePath)
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

		if surge.random {
			//https://yourbasic.org/golang/shuffle-slice-array/
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })
		}

		result = surge.execute(lines)
	}
	return
}
