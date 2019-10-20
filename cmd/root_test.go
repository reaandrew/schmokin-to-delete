package cmd_test

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/reaandrew/surge/cmd"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

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
		callback(*r)
		io.WriteString(w, "hello world\n")
	}))
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// returns ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// NOTE: there is a chance that next line won't have time to run,
			// as main() doesn't wait for this goroutine to stop. don't use
			// code with race conditions like these for production. see post
			// comments below on more discussion on how to handle this.
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

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
