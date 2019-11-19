package http

import "net/http"

type DefaultHTTPClient struct {
	client http.Client
}

func NewDefaultHTTPClient() DefaultHTTPClient {
	return DefaultHTTPClient{
		client: http.Client{},
	}
}

func (httpClient DefaultHTTPClient) Execute(request *http.Request) (*http.Response, error) {
	return httpClient.client.Do(request)
}
