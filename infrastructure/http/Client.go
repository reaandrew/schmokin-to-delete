package http

import "net/http"

type Client interface {
	Execute(request *http.Request) (*http.Response, error)
}
