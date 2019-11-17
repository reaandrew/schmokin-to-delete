package cli

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/reaandrew/surge/server"
	"github.com/reaandrew/surge/service"
	"google.golang.org/grpc"
)

func AverageFloat64(values []float64) (result float64) {
	for _, value := range values {
		result = result + value
	}
	result = result / float64(len(values))
	return
}

func Sum(values []int64) (result int64) {
	for _, value := range values {
		result += value
	}
	return
}

func Max(values []int64) (result int64) {
	for _, value := range values {
		if value > result {
			result = value
		}
	}
	return
}

func Min(values []int64) (result int64) {
	for _, value := range values {
		if result == 0 || value < result {
			result = value
		}
	}
	return
}

type SurgeCLI struct {
	workers []server.SurgeServiceClient
	//TODO: Create a configuration struct for these
	urlFilePath string
	server      bool
	serverPort  int
	serverHost  string
	processes   int
	random      bool
	workerCount int
	iterations  int
}

func (surgeCLI *SurgeCLI) Run() (result *service.SurgeResult, err error) {
	if surgeCLI.server {
		//Start the server
		server.StartServer(fmt.Sprintf("%v:%v", surgeCLI.serverHost, surgeCLI.serverPort))
	} else {
		var file *os.File
		var lines []string
		// Parsing the file is a concern of the CLI not the service.
		if surgeCLI.urlFilePath != "" {
			file, err = os.Open(surgeCLI.urlFilePath)
			if err != nil {
				return
			}
			lines = []string{}
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				lines = append(lines, line)
			}
			if err = scanner.Err(); err != nil {
				return
			}
			var wg = sync.WaitGroup{}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			responses := make(chan *server.SurgeResponse, surgeCLI.processes)
			for i := 0; i < surgeCLI.processes; i++ {
				wg.Add(1)
				portNumber := 54322 + i
				go func(port int) {
					cmd := exec.Command("./surge", "--server", "--server-host", "localhost", "--server-port", strconv.Itoa(port))
					cmd.Stdout = os.Stdout
					cmd.Start()
					//This would be better to have a synchronous wait timer
					// that would panic after a given threshold.
					// e.g. WaitFor(endpoint, 10 * time.Second)
					// and maybe tie this into the PingResponse to assert
					// on Healthy
					time.Sleep(1 * time.Second)

					conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
					if err != nil {
						log.Fatalf("did not connect: %v", err)

					}
					defer conn.Close()

					client := server.NewSurgeServiceClient(conn)
					defer client.Kill(ctx, &empty.Empty{})

					response, err := client.Run(ctx, &server.SurgeRequest{
						Iterations:  int32(surgeCLI.iterations),
						Lines:       lines,
						Random:      surgeCLI.random,
						WorkerCount: int32(surgeCLI.workerCount),
					})
					responses <- response

					if err != nil {
						panic(err)
					}
					wg.Done()
				}(portNumber)
				//Try to connect to the server and once connected add the connection to the array and start the next worker
			}

			wg.Wait()
			result = &service.SurgeResult{}
			availabilities := []float64{}
			responseTimes := []float64{}
			concurrencyRate := []float64{}
			dateReceiveRates := []float64{}
			dataSendRates := []float64{}
			failedTransactions := []int64{}
			longestTransactions := []int64{}
			shortestTransactions := []int64{}
			successfulTransactions := []int64{}
			totalBytesReceived := []int64{}
			totalBytesSent := []int64{}
			transactions := []int64{}
			transactionRates := []float64{}

			for response := range responses {
				availabilities = append(availabilities, response.Availability)
				responseTimes = append(responseTimes, response.AverageResponseTime)
				concurrencyRate = append(concurrencyRate, response.ConcurrencyRate)
				dateReceiveRates = append(dateReceiveRates, response.DataReceiveRate)
				dataSendRates = append(dataSendRates, response.DataSendRate)
				failedTransactions = append(failedTransactions, response.FailedTransactions)
				longestTransactions = append(longestTransactions, response.LongestTransaction)
				shortestTransactions = append(shortestTransactions, response.ShortestTransaction)
				successfulTransactions = append(successfulTransactions, response.SuccessfulTransactions)
				totalBytesReceived = append(totalBytesReceived, int64(response.TotalBytesReceived))
				totalBytesSent = append(totalBytesSent, int64(response.TotalBytesSent))
				transactions = append(transactions, int64(response.Transactions))
				transactionRates = append(transactionRates, response.TransactionRate)

				if len(availabilities) == surgeCLI.processes {
					close(responses)
				}
			}

			result.Availability = AverageFloat64(availabilities)
			result.AverageResponseTime = AverageFloat64(responseTimes)
			result.ConcurrencyRate = AverageFloat64(concurrencyRate)
			result.DataReceiveRate = AverageFloat64(dateReceiveRates)
			result.DataSendRate = AverageFloat64(dataSendRates)
			result.FailedTransactions = Sum(failedTransactions)
			result.LongestTransaction = Max(longestTransactions)
			result.ShortestTransaction = Min(shortestTransactions)
			result.SuccessfulTransactions = Sum(successfulTransactions)
			result.TotalBytesReceived = int(Sum(totalBytesReceived))
			result.TotalBytesSent = int(Sum(totalBytesSent))
			result.Transactions = int(Sum(transactions))
			result.TransactionRate = AverageFloat64(transactionRates)
		}
	}
	return
}