package http

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/reaandrew/surge/utils"
	"github.com/spf13/cobra"
)

type Result struct {
	TotalBytesSent     int
	TotalBytesReceived int
	Error              error
	ResponseTime       time.Duration
}

type Command struct {
	Client Client
	Timer  utils.Timer
}

func (httpCommand Command) Execute(args []string) Result {
	var verb string
	var result Result

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
			httpCommand.Timer.Start()
			response, err := httpCommand.Client.Execute(request)
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
			} else if response != nil {
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
			//Stop the timer
			result.ResponseTime = httpCommand.Timer.Stop()
			return nil
		},
	}
	command.PersistentFlags().StringVarP(&verb, "verb", "X", "GET", "")
	command.SetArgs(args)
	err := command.Execute()
	if err != nil {
		panic(err)
	}

	return result
}
