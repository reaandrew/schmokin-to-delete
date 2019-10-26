package client

import (
	"net/http"
	"sync"
)

var m = sync.Mutex{}

type FakeHTTPClient struct {
	Requests    []*http.Request
	Interceptor InterceptorFunc
}

func NewFakeHTTPClient() *FakeHTTPClient {
	return &FakeHTTPClient{
		Interceptor: func(response *http.Response) {},
	}
}

type InterceptorFunc func(response *http.Response)

func (fakeClient *FakeHTTPClient) Execute(request *http.Request) (*http.Response, error) {
	m.Lock()
	fakeClient.Requests = append(fakeClient.Requests, request)
	m.Unlock()
	response := &http.Response{
		StatusCode: http.StatusOK,
	}
	fakeClient.Interceptor(response)
	return response, nil
}
