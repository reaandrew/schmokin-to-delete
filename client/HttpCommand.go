package client

import (
	"net/http"

	"github.com/spf13/cobra"
)

type HttpCommand struct {
}

func (httpCommand HttpCommand) Execute(args []string) error {

	client := http.Client{}
	var verb string

	command := &cobra.Command{
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			request, err := http.NewRequest(verb, args[0], nil)
			if err != nil {
				return err
			}
			client.Do(request)
			return nil
		},
	}
	command.PersistentFlags().StringVarP(&verb, "verb", "X", "GET", "")
	command.SetArgs(args)
	return command.Execute()
}
