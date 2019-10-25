package cmd_test

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/reaandrew/surge/client"
	"github.com/reaandrew/surge/cmd"
	"github.com/reaandrew/surge/utils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func executeCommand(root *cobra.Command, httpClient client.HttpClient, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, httpClient, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, httpClient client.HttpClient, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	cmd.HttpClient = httpClient
	c, err = root.ExecuteC()

	fmt.Println("Buffer", buf.String())

	return c, buf.String(), err

}

func TestVisitUrlsSpecifiedInAFile(t *testing.T) {
	file := utils.CreateTestFile([]string{"http://localhost:8080/1",
		"http://localhost:8080/2",
		"http://localhost:8080/3",
		"http://localhost:8080/4",
		"http://localhost:8080/5",
	})
	defer os.Remove(file.Name())

	client := &client.FakeHTTPClient{}
	urlsVisited := []string{}

	executeCommand(cmd.RootCmd, client, "-u", file.Name())

	for _, request := range client.Requests {
		urlsVisited = append(urlsVisited, request.RequestURI)
	}

	assert.Equal(t, len(urlsVisited), 5)
}

func TestSupportForVerbPut(t *testing.T) {
	file := utils.CreateTestFile([]string{
		"-X PUT http://localhost:8080/1",
		"http://localhost:8080/2 -X PUT",
	})
	defer os.Remove(file.Name())

	methods := []string{}
	client := &client.FakeHTTPClient{}

	executeCommand(cmd.RootCmd, client, "-u", file.Name())

	for _, request := range client.Requests {
		methods = append(methods, request.Method)
	}

	assert.Equal(t, methods, []string{"PUT", "PUT"})
}

func TestSupportForRandomOrder(t *testing.T) {
	urls := func() []string {
		returnUrls := []string{}
		for i := 0; i < 10; i++ {
			returnUrls = append(returnUrls, "http://localhost:8080/"+strconv.Itoa(i))
		}
		return returnUrls
	}()
	file := utils.CreateTestFile(urls)
	defer os.Remove(file.Name())

	client := &client.FakeHTTPClient{}

	output, err := executeCommand(cmd.RootCmd, client, "-u", file.Name(), "-r")
	assert.Nil(t, err, output)

	urlsVisited := []string{}
	for _, request := range client.Requests {
		urlsVisited = append(urlsVisited, request.RequestURI)
	}

	urlPaths := utils.MapStrings(urls, func(value string) string {
		items := strings.Split(value, "/")
		return "/" + items[len(items)-1]
	})
	assert.NotEqual(t, urlsVisited, urlPaths)
}

func TestSupportForConcurrentWorkers(t *testing.T) {
	file := utils.CreateTestFile([]string{
		"http://localhost:8080/1",
	})
	defer os.Remove(file.Name())

	var concurrentWorkerCount = 5
	client := &client.FakeHTTPClient{}

	output, err := executeCommand(cmd.RootCmd, client, "-u", file.Name(), "-c", strconv.Itoa(concurrentWorkerCount))
	assert.Nil(t, err, output)
	assert.Equal(t, len(client.Requests), concurrentWorkerCount)
}

func TestSupportForNumberOfIterations(t *testing.T) {
	file := utils.CreateTestFile([]string{
		"http://localhost:8080/2",
	})
	defer os.Remove(file.Name())

	var concurrentWorkerCount = 1
	var iterationCount = 5
	client := &client.FakeHTTPClient{}

	output, err := executeCommand(cmd.RootCmd, client, "-u", file.Name(),
		"-n", strconv.Itoa(iterationCount),
		"-c", strconv.Itoa(concurrentWorkerCount))
	assert.Nil(t, err, output)
	assert.Equal(t, len(client.Requests), iterationCount)
}

func TestSupportForNumberOfIterationsWithConcurrentWorkers(t *testing.T) {
	file := utils.CreateTestFile([]string{
		"http://localhost:8080/1",
	})
	defer os.Remove(file.Name())

	var concurrentWorkerCount = 5
	var iterationCount = 5
	client := &client.FakeHTTPClient{}

	output, err := executeCommand(cmd.RootCmd, client, "-u", file.Name(),
		"-n", strconv.Itoa(iterationCount),
		"-c", strconv.Itoa(concurrentWorkerCount))
	assert.Nil(t, err, output)
	assert.Equal(t, len(client.Requests), iterationCount*concurrentWorkerCount)
}

func TestOutputsNumberOfTransactions(t *testing.T) {
	file := utils.CreateTestFile([]string{
		"http://localhost:8080/1",
	})
	defer os.Remove(file.Name())

	client := &client.FakeHTTPClient{}

	output, err := executeCommand(cmd.RootCmd, client, "-u", file.Name(), "-n", "1", "-c", "1")

	assert.Nil(t, err)
	assert.Contains(t, output, "Transactions: 1\n")
}
