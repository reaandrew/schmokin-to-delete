package server

import (
	context "context"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/rcrowley/go-metrics"
	"github.com/reaandrew/surge/core"
	"github.com/reaandrew/surge/utils"
	grpc "google.golang.org/grpc"
)

const port = ":50051"

type surgeRemoteService struct {
	random      bool
	workerCount int
	iterations  int
	processes   int
	httpClient  core.HttpClient
	timer       utils.Timer
	lock        sync.Mutex
	waitGroup   sync.WaitGroup
	server      bool
	serverPort  int
	serverHost  string
	//TODO: Create a stats struct for these
	transactions           int32
	errors                 int32
	totalBytesSent         int32
	totalBytesReceived     int32
	responseTime           metrics.Histogram
	transactionRate        metrics.Meter
	concurrencyCounter     metrics.Counter
	concurrencyRate        metrics.Histogram
	dataSendRate           metrics.Meter
	dataReceiveRate        metrics.Meter
	successfulTransactions int32
}

func (s *surgeRemoteService) Run(ctx context.Context, in *SurgeRequest) (*SurgeResponse, error) {
	return s.execute(in.Lines), nil
}

func (surge *surgeRemoteService) worker(linesValue []string) {
	for i := 0; i < len(linesValue) || (surge.iterations > 0 && i < surge.iterations); i++ {
		line := linesValue[i%len(linesValue)]
		var command = HttpCommand{
			client: surge.httpClient,
			timer:  surge.timer,
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

func (surge *surgeRemoteService) execute(lines []string) *SurgeResponse {
	for i := 0; i < surge.workerCount; i++ {
		surge.timer.Start()
		surge.waitGroup.Add(1)
		go surge.worker(lines)
	}
	surge.waitGroup.Wait()
	result := &SurgeResponse{
		Transactions:           surge.transactions,
		ElapsedTime:            int64(surge.timer.Stop()),
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

func StartServer() *grpc.Server {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}

	server := grpc.NewServer()
	RegisterSurgeServiceServer(server, &surgeRemoteService{})

	if err := server.Serve(lis); err != nil {
		log.Fatal(errors.Wrap(err, "Failed to start server!"))
	}

	return server
}
