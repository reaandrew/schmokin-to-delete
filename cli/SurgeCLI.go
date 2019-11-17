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
		fmt.Println("Starting Server!")
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

					fmt.Println("About to start remote worker")
					response, err := client.Run(ctx, &server.SurgeRequest{
						Iterations:  int32(surgeCLI.iterations),
						Lines:       lines,
						Random:      surgeCLI.random,
						WorkerCount: int32(surgeCLI.workerCount),
					})
					fmt.Println("Finished", err)

					if err != nil {
						panic(err)
					} else {
						fmt.Println("Response", response)
					}
					wg.Done()
				}(portNumber)
				//Try to connect to the server and once connected add the connection to the array and start the next worker
			}

			wg.Wait()
		}
	}
	return &service.SurgeResult{}, nil
}
