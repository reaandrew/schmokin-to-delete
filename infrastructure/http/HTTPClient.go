package http

import "net/http"

type HTTPClient interface {
	Execute(request *http.Request) (*http.Response, error)
}
