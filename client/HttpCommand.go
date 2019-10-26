package client

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

type HttpCommand struct {
	client HttpClient
}

func (httpCommand HttpCommand) Execute(args []string) error {
	var verb string
	var returnError error

	command := &cobra.Command{
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			request, err := http.NewRequest(verb, args[0], nil)
			if err != nil {
				return err
			}
			response, err := httpCommand.client.Execute(request)
			if err != nil {
				return err
			}
			if response.StatusCode != 200 {
				returnError = errors.New("Error " + strconv.Itoa(response.StatusCode))
			}
			return nil
		},
	}
	command.PersistentFlags().StringVarP(&verb, "verb", "X", "GET", "")
	command.SetArgs(args)
	command.Execute()

	return returnError
}
