package client

import "net/http"

type DefaultHttpClient struct {
	client http.Client
}

func NewDefaultHttpClient() DefaultHttpClient {
	return DefaultHttpClient{
		client: http.Client{},
	}
}

func (httpClient DefaultHttpClient) Execute(request *http.Request) {
	httpClient.client.Do(request)
}
