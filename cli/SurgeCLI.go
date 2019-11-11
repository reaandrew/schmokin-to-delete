package cli

import (
	"bufio"
	"math/rand"
	"os"
	"time"

	"github.com/reaandrew/surge/server"
)

type SurgeCLI struct {
	//TODO: Create a configuration struct for these
	urlFilePath string
	server      bool
	serverPort  int
	serverHost  string
	workers     []server.SurgeServiceClient
	processes   int
}

func (SurgeCLI *SurgeCLI) Run() (result Result, err error) {
	var file *os.File
	if SurgeCLI.urlFilePath != "" {
		file, err = os.Open(SurgeCLI.urlFilePath)
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

		if SurgeCLI.random {
			//https://yourbasic.org/golang/shuffle-slice-array/
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })
		}

		if SurgeCLI.server {
			//Start the server
			server.StartServer()
		} else {
			for 
		}

	}
	return
}
