package service

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
	schmokinHTTP "github.com/reaandrew/schmokin/infrastructure/http"
	"github.com/reaandrew/schmokin/utils"
)

type SchmokinService struct {
	random      bool
	workerCount int
	iterations  int
	httpClient  schmokinHTTP.Client
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

func (schmokin *SchmokinService) worker(linesValue []string) {
	for i := 0; i < len(linesValue) || (schmokin.iterations > 0 && i < schmokin.iterations); i++ {
		line := linesValue[i%len(linesValue)]
		var command = schmokinHTTP.Command{
			Client: schmokin.httpClient,
			Timer:  schmokin.timer,
		}
		var args = strings.Fields(line)
		schmokin.concurrencyCounter.Inc(1)
		result := command.Execute(args)
		schmokin.concurrencyCounter.Dec(1)
		schmokin.concurrencyRate.Update(schmokin.concurrencyCounter.Count())
		schmokin.lock.Lock()
		if result.Error != nil {
			schmokin.errors++
		} else {
			schmokin.successfulTransactions++
		}
		schmokin.transactions++
		schmokin.totalBytesSent += result.TotalBytesSent
		schmokin.totalBytesReceived += result.TotalBytesReceived
		schmokin.responseTime.Update(int64(result.ResponseTime))
		schmokin.dataSendRate.Mark(int64(result.TotalBytesSent))
		schmokin.dataReceiveRate.Mark(int64(result.TotalBytesReceived))
		schmokin.transactionRate.Mark(1)
		schmokin.lock.Unlock()
		if i > 0 && i == schmokin.iterations-1 {
			break
		}
	}
	schmokin.waitGroup.Done()
}

func (schmokin *SchmokinService) Execute(lines []string) SchmokinResult {
	timer := schmokin.timer.Start()
	if schmokin.random {
		//https://yourbasic.org/golang/shuffle-slice-array/
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })
	}
	for i := 0; i < schmokin.workerCount; i++ {
		schmokin.waitGroup.Add(1)
		go schmokin.worker(lines)
	}
	schmokin.waitGroup.Wait()
	result := SchmokinResult{
		Transactions:           schmokin.transactions,
		ElapsedTime:            timer.Stop(),
		TotalBytesSent:         schmokin.totalBytesSent,
		TotalBytesReceived:     schmokin.totalBytesReceived,
		AverageResponseTime:    schmokin.responseTime.Mean(),
		TransactionRate:        schmokin.transactionRate.RateMean(),
		ConcurrencyRate:        schmokin.concurrencyRate.Mean(),
		DataSendRate:           schmokin.dataSendRate.RateMean(),
		DataReceiveRate:        schmokin.dataReceiveRate.RateMean(),
		SuccessfulTransactions: int64(schmokin.successfulTransactions),
		FailedTransactions:     int64(schmokin.errors),
		LongestTransaction:     schmokin.responseTime.Max(),
		ShortestTransaction:    schmokin.responseTime.Min(),
	}
	if schmokin.errors == 0 {
		result.Availability = 1
	} else {
		availability := float64(schmokin.errors) / float64(schmokin.transactions)
		if availability < 1 {
			result.Availability = float64(1) - availability
		} else {
			if schmokin.errors == schmokin.transactions {
				result.Availability = 0
			} else {
				result.Availability = availability
			}
		}
	}
	return result
}
