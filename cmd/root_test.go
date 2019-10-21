package cmd_test

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/gorilla/mux"
	"github.com/reaandrew/surge/cmd"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var m = sync.Mutex{}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	c, err = root.ExecuteC()

	return c, buf.String(), err

}

func startHTTPServer(callback func(r http.Request)) *http.Server {
	r := mux.NewRouter()
	r.Handle("/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Lock()
		callback(*r)
		m.Unlock()
		io.WriteString(w, "hello world\n")
	}))
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go srv.ListenAndServe()

	// returning reference so caller can call Shutdown()
	return srv
}

func TestVisitUrlsSpecifiedInAFile(t *testing.T) {
	file := CreateTestFile([]string{"http://localhost:8080/1",
		"http://localhost:8080/2",
		"http://localhost:8080/3",
		"http://localhost:8080/4",
		"http://localhost:8080/5",
	})
	defer os.Remove(file.Name())

	urlsVisited := []string{}
	srv := startHTTPServer(func(r http.Request) {
		urlsVisited = append(urlsVisited, r.RequestURI)
	})
	defer srv.Shutdown(context.TODO())

	executeCommand(cmd.RootCmd, "-u", file.Name())

	assert.Equal(t, len(urlsVisited), 5)
}

func TestSupportForVerbPut(t *testing.T) {
	file := CreateTestFile([]string{
		"-X PUT http://localhost:8080/1",
		"http://localhost:8080/2 -X PUT",
	})
	defer os.Remove(file.Name())

	methods := []string{}
	srv := startHTTPServer(func(r http.Request) {
		methods = append(methods, r.Method)
	})
	defer srv.Shutdown(context.TODO())

	executeCommand(cmd.RootCmd, "-u", file.Name())

	assert.Equal(t, methods, []string{"PUT", "PUT"})
}

func CreateTestFile(lines []string) *os.File {
	fileContents := strings.Join(lines, "\n")
	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(file.Name(), []byte(fileContents), os.ModePerm)
	return file
}

func MapStrings(array []string, delegate func(value string) string) (values []string) {
	for _, value := range array {
		values = append(values, delegate(value))
	}
	return
}

func TestSupportForRandomOrder(t *testing.T) {
	urls := func() []string {
		returnUrls := []string{}
		for i := 0; i < 10; i++ {
			returnUrls = append(returnUrls, "http://localhost:8080/"+strconv.Itoa(i))
		}
		return returnUrls
	}()
	file := CreateTestFile(urls)
	defer os.Remove(file.Name())

	urlsVisited := []string{}
	srv := startHTTPServer(func(r http.Request) {
		urlsVisited = append(urlsVisited, r.RequestURI)
	})
	defer srv.Shutdown(context.TODO())

	output, err := executeCommand(cmd.RootCmd, "-u", file.Name(), "-r")
	assert.Nil(t, err, output)

	urlPaths := MapStrings(urls, func(value string) string {
		items := strings.Split(value, "/")
		return "/" + items[len(items)-1]
	})
	assert.NotEqual(t, urlsVisited, urlPaths)
}

func TestSupportForConcurrentWorkers(t *testing.T) {
	file := CreateTestFile([]string{
		"http://localhost:8080/1",
	})
	defer os.Remove(file.Name())

	var concurrentWorkerCount = 5
	var count int
	srv := startHTTPServer(func(r http.Request) {
		count++
	})
	defer srv.Shutdown(context.TODO())

	output, err := executeCommand(cmd.RootCmd, "-u", file.Name(), "-c", strconv.Itoa(concurrentWorkerCount))
	assert.Nil(t, err, output)
	assert.Equal(t, count, concurrentWorkerCount)
}
