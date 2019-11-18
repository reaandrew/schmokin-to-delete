package cli

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/reaandrew/surge/server"
	"github.com/reaandrew/surge/service"
	"github.com/reaandrew/surge/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
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

func ReadFileToLines(path string) (lines []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	lines = []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	err = scanner.Err()
	return
}

type SurgeServiceClientConnection struct {
	Client     server.SurgeServiceClient
	Connection *grpc.ClientConn
}

type SurgeCLI struct {
	workers []SurgeServiceClientConnection
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

func (surgeCLI *SurgeCLI) StartServer(port int) SurgeServiceClientConnection {
	cmd := exec.Command("./surge", "--server", "--server-host", "localhost", "--server-port", strconv.Itoa(port))
	cmd.Stdout = os.Stdout
	cmd.Start()
	//This would be better to have a synchronous wait timer
	// that would panic after a given threshold.
	// e.g. WaitFor(endpoint, 10 * time.Second)
	// and maybe tie this into the PingResponse to assert
	// on Healthy

	var conn *grpc.ClientConn
	var err error

	utils.WaitUtil{
		Timeout: 1 * time.Second,
		Backoff: 10 * time.Millisecond,
	}.Wait(func() bool {
		conn, err = grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
		return err == nil
	})

	utils.WaitUtil{
		Timeout: 10 * time.Second,
		Backoff: 250 * time.Millisecond,
	}.Wait(func() bool {
		return conn.GetState() == connectivity.Ready
	})

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := server.NewSurgeServiceClient(conn)

	return SurgeServiceClientConnection{
		Connection: conn,
		Client:     client,
	}
}

func (surgeCLI *SurgeCLI) RunServer() (result *service.SurgeResult, err error) {
	server.StartServer(fmt.Sprintf("%v:%v", surgeCLI.serverHost, surgeCLI.serverPort))
	return
}

func (surgeCLI *SurgeCLI) RunController() (result *service.SurgeResult, err error) {
	var lines []string

	if surgeCLI.urlFilePath == "" {
		panic("No URL file supplied")
	}

	lines, err = ReadFileToLines(surgeCLI.urlFilePath)
	if err != nil {
		panic(err)
	}
	var wg = sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	responses := make(chan *server.SurgeResponse, surgeCLI.processes)

	fmt.Println("Starting the worker processes...")

	for i := 0; i < surgeCLI.processes; i++ {
		portNumber := 54322 + i
		connection := surgeCLI.StartServer(portNumber)
		defer func() {
			connection.Client.Kill(ctx, &empty.Empty{})
			connection.Connection.Close()
		}()
		surgeCLI.workers = append(surgeCLI.workers, connection)
	}

	fmt.Println("Surging...")

	for _, connection := range surgeCLI.workers {
		wg.Add(1)
		go func(connection SurgeServiceClientConnection) {
			response, err := connection.Client.Run(ctx, &server.SurgeRequest{
				Iterations:  int32(surgeCLI.iterations),
				Lines:       lines,
				Random:      surgeCLI.random,
				WorkerCount: int32(surgeCLI.workerCount),
			})
			responses <- response

			if err != nil {
				fmt.Println(connection.Connection.GetState())
				panic(err)
			}
			wg.Done()
		}(connection)
	}

	wg.Wait()
	fmt.Println("Stopping the worker processes...")
	for _, connection := range surgeCLI.workers {
		wg.Add(1)
		go func(connection SurgeServiceClientConnection) {
			utils.WaitUtil{
				Timeout: 10 * time.Second,
				Backoff: 250 * time.Millisecond,
			}.Wait(func() bool {
				_, err := connection.Client.Kill(ctx, &empty.Empty{})
				return err == nil || strings.Contains(err.Error(), "Unavailable")
			})
			wg.Done()
		}(connection)
	}
	wg.Wait()

	close(responses)
	result = surgeCLI.MergeResponses(responses)
	return
}

func (surgeCLI *SurgeCLI) Run() (result *service.SurgeResult, err error) {
	if surgeCLI.server {
		return surgeCLI.RunServer()
	}

	return surgeCLI.RunController()
}

func (surgeCLI SurgeCLI) MergeResponses(responses chan *server.SurgeResponse) (result *service.SurgeResult) {
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
	return
}
