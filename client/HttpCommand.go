package client

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/reaandrew/surge/utils"
	"github.com/spf13/cobra"
)

type HttpResult struct {
	TotalBytesSent     int
	TotalBytesReceived int
	Error              error
	ResponseTime       time.Duration
}

type HttpCommand struct {
	client HttpClient
	timer  utils.Timer
}

func (httpCommand HttpCommand) Execute(args []string) HttpResult {
	var verb string
	var result HttpResult

	command := &cobra.Command{
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			request, err := http.NewRequest(verb, args[0], nil)
			if err != nil {
				return err
			}
			requestBytes, err := httputil.DumpRequestOut(request, true)
			if err != nil {
				return err
			}
			result.TotalBytesSent = len(requestBytes)
			//When using the TRACE utility for HTTP with golang
			// we can still use the Timer interface
			//Start the timer
			httpCommand.timer.Start()
			response, err := httpCommand.client.Execute(request)
			if err != nil {
				result.Error = err
			} else {
				if response.Body != nil {
					defer response.Body.Close()
					io.Copy(ioutil.Discard, response.Body)
				}
				responseBytes, err := httputil.DumpResponse(response, true)
				result.TotalBytesReceived = len(responseBytes)
				if err != nil {
					result.Error = err
				}
				if response.StatusCode >= 400 {
					result.Error = errors.New("Error " + strconv.Itoa(response.StatusCode))
				}
			}
			//Stop the timer
			result.ResponseTime = httpCommand.timer.Stop()
			return nil
		},
	}
	command.PersistentFlags().StringVarP(&verb, "verb", "X", "GET", "")
	command.SetArgs(args)
	command.Execute()

	return result
}
