package client

import "net/http"

type HttpClient interface {
	Execute(request *http.Request) (*http.Response, error)
}
