package client

import (
	"bufio"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type Surge struct {
	UrlFilePath string
}

func (surge Surge) Run() error {
	if surge.UrlFilePath != "" {
		file, err := os.Open(surge.UrlFilePath)
		if err != nil {
			return err
		}
		client := http.Client{}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

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
			command.SetArgs(strings.Fields(line))
			command.Execute()

		}

		if err := scanner.Err(); err != nil {
			return err
		}
	}
	return nil
}
