package http

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/reaandrew/schmokin/utils"
	"github.com/urfave/cli/v2"
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
	verb   string
}

func (httpCommand Command) run(args []string) (result Result) {
	var verb = httpCommand.verb

	request, err := http.NewRequest(verb, args[0], nil)
	if err != nil {
		result.Error = err
		return
	}
	// When using the TRACE utility for HTTP with golang
	// we can still use the Timer interface

	// Start the timer
	timer := httpCommand.Timer.Start()
	response, err := httpCommand.Client.Execute(request)
	if err != nil {
		result.Error = err
		return
	}
	requestBytes, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		result.Error = err
		return
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
			result.Error = err
		}
		result.TotalBytesReceived = len(responseBytes)
		if err != nil {
			result.Error = err
		}
		if response.StatusCode >= 400 {
			result.Error = errors.New("Error " + strconv.Itoa(response.StatusCode))
		}
	}
	// Stop the timer
	result.ResponseTime = timer.Stop()
	return
}

func (httpCommand Command) Execute(args []string) Result {
	var result Result

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "verb",
				Value:       "POST",
				Aliases:     []string{"X"},
				Usage:       "verb",
				Destination: &httpCommand.verb,
			},
			&cli.StringSliceFlag{
				Name:    "header",
				Usage:   "header",
				Aliases: []string{"H"},
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println(c.StringSlice("header"))
			result = httpCommand.run(args)
			return nil
		},
	}
	app.Run(args)

	return result
}
