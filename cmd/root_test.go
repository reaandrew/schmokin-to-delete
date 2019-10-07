package cmd_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/reaandrew/surge/cmd"
	"github.com/spf13/cobra"
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
	srv := &http.Server{Addr: ":8080"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		callback(*r)
		io.WriteString(w, "hello world\n")
	})

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
	fileContents := `http://localhost:8080/1
http://localhost:8080/2
http://localhost:8080/3
http://localhost:8080/4
http://localhost:8080/5`

	//Write the urls to a file
	//Pass the urls file path in as a -u flag
	file, err := ioutil.TempFile(os.TempDir(), "prefix")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	ioutil.WriteFile(file.Name(), []byte(fileContents), os.ModePerm)

	urlsVisited := []string{}
	srv := startHTTPServer(func(r http.Request) {
		urlsVisited = append(urlsVisited, r.RequestURI)
	})

	output, err := executeCommand(cmd.RootCmd, "-u", file.Name())
	fmt.Println(urlsVisited, output)

	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}
