package client

import (
	"bufio"
	"net/http"
	"os"
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
			request, err := http.NewRequest("GET", scanner.Text(), nil)
			if err != nil {
				return err
			}
			client.Do(request)
		}

		if err := scanner.Err(); err != nil {
			return err
		}
	}
	return nil
}
