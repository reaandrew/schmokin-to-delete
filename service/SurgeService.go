package service

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
	surgeHTTP "github.com/reaandrew/surge/infrastructure/http"
	"github.com/reaandrew/surge/utils"
)

type SurgeService struct {
	random      bool
	workerCount int
	iterations  int
	httpClient  surgeHTTP.HttpClient
	timer       utils.Timer
	lock        sync.Mutex
	waitGroup   sync.WaitGroup
	//TODO: Create a stats struct for these
	transactions           int
	errors                 int
	totalBytesSent         int
	totalBytesReceived     int
	responseTime           metrics.Histogram
	transactionRate        metrics.Meter
	concurrencyCounter     metrics.Counter
	concurrencyRate        metrics.Histogram
	dataSendRate           metrics.Meter
	dataReceiveRate        metrics.Meter
	successfulTransactions int
}

func (surge *SurgeService) worker(linesValue []string) {
	for i := 0; i < len(linesValue) || (surge.iterations > 0 && i < surge.iterations); i++ {
		line := linesValue[i%len(linesValue)]
		var command = surgeHTTP.HttpCommand{
			Client: surge.httpClient,
			Timer:  surge.timer,
		}
		var args = strings.Fields(line)
		surge.concurrencyCounter.Inc(1)
		result := command.Execute(args)
		surge.concurrencyCounter.Dec(1)
		surge.concurrencyRate.Update(surge.concurrencyCounter.Count())
		surge.lock.Lock()
		if result.Error != nil {
			surge.errors++
		} else {
			surge.successfulTransactions++
		}
		surge.transactions++
		surge.totalBytesSent += result.TotalBytesSent
		surge.totalBytesReceived += result.TotalBytesReceived
		surge.responseTime.Update(int64(result.ResponseTime))
		surge.dataSendRate.Mark(int64(result.TotalBytesSent))
		surge.dataReceiveRate.Mark(int64(result.TotalBytesReceived))
		surge.transactionRate.Mark(1)
		surge.lock.Unlock()
		if i > 0 && i == surge.iterations-1 {
			break
		}
	}
	surge.waitGroup.Done()
}

func (surge *SurgeService) Execute(lines []string) SurgeResult {
	if surge.random {
		//https://yourbasic.org/golang/shuffle-slice-array/
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })
	}

	for i := 0; i < surge.workerCount; i++ {
		surge.timer.Start()
		surge.waitGroup.Add(1)
		go surge.worker(lines)
	}
	surge.waitGroup.Wait()
	result := SurgeResult{
		Transactions:           surge.transactions,
		ElapsedTime:            surge.timer.Stop(),
		TotalBytesSent:         surge.totalBytesSent,
		TotalBytesReceived:     surge.totalBytesReceived,
		AverageResponseTime:    surge.responseTime.Mean(),
		TransactionRate:        surge.transactionRate.RateMean(),
		ConcurrencyRate:        float64(surge.concurrencyRate.Mean()),
		DataSendRate:           surge.dataSendRate.RateMean(),
		DataReceiveRate:        surge.dataReceiveRate.RateMean(),
		SuccessfulTransactions: int64(surge.successfulTransactions),
		FailedTransactions:     int64(surge.errors),
		LongestTransaction:     surge.responseTime.Max(),
		ShortestTransaction:    surge.responseTime.Min(),
	}
	if surge.errors == 0 {
		result.Availability = 1
	} else {
		availability := float64(surge.errors) / float64(surge.transactions)
		if availability < 1 {
			result.Availability = float64(1 - availability)
		} else {
			if surge.errors == surge.transactions {
				result.Availability = 0
			} else {
				result.Availability = availability
			}
		}
	}
	return result
}
