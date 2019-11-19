package http

import "net/http"

type DefaultClient struct {
	client http.Client
}

func NewDefaultClient() DefaultClient {
	return DefaultClient{
		client: http.Client{},
	}
}

func (httpClient DefaultClient) Execute(request *http.Request) (*http.Response, error) {
	return httpClient.client.Do(request)
}
