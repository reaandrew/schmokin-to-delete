package cli

import (
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
	var err error
	cmd := exec.Command("./surge", "--server", "--server-host", "localhost", "--server-port", strconv.Itoa(port))
	cmd.Stdout = os.Stdout
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	//This would be better to have a synchronous wait timer
	// that would panic after a given threshold.
	// e.g. WaitFor(endpoint, 10 * time.Second)
	// and maybe tie this into the PingResponse to assert
	// on Healthy

	var conn *grpc.ClientConn

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

func (surgeCLI *SurgeCLI) StartWorkerProcesses() {
	var wg = sync.WaitGroup{}
	var syncLock = sync.Mutex{}
	for i := 0; i < surgeCLI.processes; i++ {
		wg.Add(1)
		//This needs to use some sort of freeport package to find any port which is going
		// rather than a known range, which currently is very limiting
		portNumber := 54322 + i
		go func(port int) {
			connection := surgeCLI.StartServer(portNumber)
			syncLock.Lock()
			surgeCLI.workers = append(surgeCLI.workers, connection)
			syncLock.Unlock()
			wg.Done()
		}(portNumber)
	}
	wg.Wait()
}

func (surgeCLI *SurgeCLI) ExecuteWorkerProcesses(ctx context.Context, lines []string) (responses []*server.SurgeResponse) {
	var wg = sync.WaitGroup{}
	var lock = sync.Mutex{}
	for _, connection := range surgeCLI.workers {
		wg.Add(1)
		go func(connection SurgeServiceClientConnection) {
			response, err := connection.Client.Run(ctx, &server.SurgeRequest{
				Iterations:  int32(surgeCLI.iterations),
				Lines:       lines,
				Random:      surgeCLI.random,
				WorkerCount: int32(surgeCLI.workerCount),
			})
			lock.Lock()
			responses = append(responses, response)
			lock.Unlock()
			if err != nil {
				fmt.Println(connection.Connection.GetState())
				panic(err)
			}
			wg.Done()
		}(connection)
	}

	wg.Wait()
	return
}

func (surgeCLI *SurgeCLI) StopWorkerProcesses(ctx context.Context) {
	var wg = sync.WaitGroup{}
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
}

func (surgeCLI *SurgeCLI) RunController() (result *service.SurgeResult, err error) {
	var lines []string

	if surgeCLI.urlFilePath == "" {
		panic("No URL file supplied")
	}

	lines, err = utils.ReadFileToLines(surgeCLI.urlFilePath)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Starting the worker processes...")
	surgeCLI.StartWorkerProcesses()

	fmt.Println("Surging...")
	responses := surgeCLI.ExecuteWorkerProcesses(ctx, lines)

	fmt.Println("Stopping the worker processes...")
	surgeCLI.StopWorkerProcesses(ctx)

	result = server.MergeResponses(responses)
	return
}

func (surgeCLI *SurgeCLI) Run() (result *service.SurgeResult, err error) {
	if surgeCLI.server {
		return surgeCLI.RunServer()
	}

	return surgeCLI.RunController()
}
