package client

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/reaandrew/surge/core"
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
	client core.HttpClient
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
			//When using the TRACE utility for HTTP with golang
			// we can still use the Timer interface
			//Start the timer
			httpCommand.timer.Start()
			response, err := httpCommand.client.Execute(request)
			if err != nil {
				result.Error = err
				return nil
			}
			requestBytes, err := httputil.DumpRequestOut(request, true)
			if err != nil {
				result.Error = err
				return nil
			}
			result.TotalBytesSent = len(requestBytes)
			if err != nil {
				result.Error = err
			} else {
				if response != nil {
					if response.Body != nil {
						defer response.Body.Close()
					}
					responseBytes, err := httputil.DumpResponse(response, true)
					if err != nil {
						panic(err)
					}
					result.TotalBytesReceived = len(responseBytes)
					if err != nil {
						result.Error = err
					}
					if response.StatusCode >= 400 {
						result.Error = errors.New("Error " + strconv.Itoa(response.StatusCode))
					}
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
