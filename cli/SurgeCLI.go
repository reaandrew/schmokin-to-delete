package cli

import (
	"bufio"
	"os"
	"os/exec"
	"strconv"

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
		for i := 0; i < surgeCLI.processes; i++ {
			portNumber := 54322 + i
			cmd := exec.Command("surge", "--server", "--server-host", "localhost", "--server-port", strconv.Itoa(portNumber))
			cmd.Start()
			//Try to connect to the server and once connected add the connection to the array and start the next worker
		}
	}
	return
}
