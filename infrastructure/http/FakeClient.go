package http

import (
	"net/http"
	"sync"
)

var m = sync.Mutex{}

type FakeClient struct {
	Requests    []*http.Request
	Interceptor InterceptorFunc
}

func NewFakeClient() *FakeClient {
	return &FakeClient{
		Interceptor: func(response *http.Response) {},
	}
}

type InterceptorFunc func(response *http.Response)

func (fakeClient *FakeClient) Execute(request *http.Request) (*http.Response, error) {
	m.Lock()
	fakeClient.Requests = append(fakeClient.Requests, request)
	m.Unlock()
	response := &http.Response{
		StatusCode: http.StatusOK,
	}
	fakeClient.Interceptor(response)
	return response, nil
}
