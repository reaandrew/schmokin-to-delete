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
	"github.com/reaandrew/schmokin/server"
	"github.com/reaandrew/schmokin/service"
	"github.com/reaandrew/schmokin/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type SchmokinServiceClientConnection struct {
	Client     server.SchmokinServiceClient
	Connection *grpc.ClientConn
}

type SchmokinCLI struct {
	workers []SchmokinServiceClientConnection
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

const SchmokinPathVar = "SCHMOKIN_PATH"

func (schmokinCLI *SchmokinCLI) StartServer(port int) SchmokinServiceClientConnection {
	var err error

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	if os.Getenv(SchmokinPathVar) != "" {
		ex = os.Getenv(SchmokinPathVar)
	}

	cmd := exec.Command(ex, "--server", "--server-host", "localhost", "--server-port", strconv.Itoa(port))
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	// This would be better to have a synchronous wait timer
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

	client := server.NewSchmokinServiceClient(conn)

	return SchmokinServiceClientConnection{
		Connection: conn,
		Client:     client,
	}
}

func (schmokinCLI *SchmokinCLI) RunServer() (result *service.SchmokinResult, err error) {
	log.Println(fmt.Sprintf("Starting server %s %d", schmokinCLI.serverHost, schmokinCLI.serverPort))
	server.StartServer(fmt.Sprintf("%v:%v", schmokinCLI.serverHost, schmokinCLI.serverPort))
	return &service.SchmokinResult{}, nil
}

func (schmokinCLI *SchmokinCLI) StartWorkerProcesses() {
	var wg = sync.WaitGroup{}
	var syncLock = sync.Mutex{}
	for i := 0; i < schmokinCLI.processes; i++ {
		wg.Add(1)
		// This needs to use some sort of freeport package to find any port which is going
		// rather than a known range, which currently is very limiting
		portNumber := 54322 + i
		go func(port int) {
			connection := schmokinCLI.StartServer(port)
			syncLock.Lock()
			schmokinCLI.workers = append(schmokinCLI.workers, connection)
			syncLock.Unlock()
			wg.Done()
		}(portNumber)
	}
	wg.Wait()
}

func (schmokinCLI *SchmokinCLI) ExecuteWorkerProcesses(ctx context.Context, lines []string) (responses []*server.SchmokinResponse) {
	var wg = sync.WaitGroup{}
	var lock = sync.Mutex{}
	for _, connection := range schmokinCLI.workers {
		wg.Add(1)
		go func(connection SchmokinServiceClientConnection) {
			response, err := connection.Client.Run(ctx, &server.SchmokinRequest{
				Iterations:  int32(schmokinCLI.iterations),
				Lines:       lines,
				Random:      schmokinCLI.random,
				WorkerCount: int32(schmokinCLI.workerCount),
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

func (schmokinCLI *SchmokinCLI) StopWorkerProcesses(ctx context.Context) {
	var wg = sync.WaitGroup{}
	for _, connection := range schmokinCLI.workers {
		wg.Add(1)
		go func(connection SchmokinServiceClientConnection) {
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

func (schmokinCLI *SchmokinCLI) RunController() (result *service.SchmokinResult, err error) {
	var lines []string

	if schmokinCLI.urlFilePath == "" {
		panic("No URL file supplied")
	}

	lines, err = utils.ReadFileToLines(schmokinCLI.urlFilePath)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Starting the worker processes...")
	schmokinCLI.StartWorkerProcesses()

	fmt.Println("Surging...")
	responses := schmokinCLI.ExecuteWorkerProcesses(ctx, lines)

	fmt.Println("Stopping the worker processes...")
	schmokinCLI.StopWorkerProcesses(ctx)

	result = server.MergeResponses(responses)
	return
}

func (schmokinCLI *SchmokinCLI) Run() (result *service.SchmokinResult, err error) {
	if schmokinCLI.server {
		return schmokinCLI.RunServer()
	}

	return schmokinCLI.RunController()
}
