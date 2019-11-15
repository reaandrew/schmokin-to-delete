package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/reaandrew/surge/server"
	"github.com/reaandrew/surge/service"
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
		server.StartServer()
	} else {
		var file *os.File
		// Parsing the file is a concern of the CLI not the service.
		if surgeCLI.urlFilePath != "" {
			file, err = os.Open(surgeCLI.urlFilePath)
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
		}
		var wg = sync.WaitGroup{}
		for i := 0; i < surgeCLI.processes; i++ {
			wg.Add(1)
			portNumber := 54322 + i
			go func(port int) {
				cmd := exec.Command("./surge", "--server", "--server-host", "localhost", "--server-port", strconv.Itoa(port))
				cmd.Stdout = os.Stdout
				fmt.Println("Starting", strconv.Itoa(port))
				err := cmd.Run()
				fmt.Println("Finished", err)
				wg.Done()
			}(portNumber)
			//Try to connect to the server and once connected add the connection to the array and start the next worker
		}
		wg.Wait()
	}
	return &service.SurgeResult{}, nil
}
