package client

import (
	"net/http"
	"sync"
)

var m = sync.Mutex{}

type FakeHTTPClient struct {
	Requests []*http.Request
}

func (fakeClient *FakeHTTPClient) Execute(request *http.Request) {
	m.Lock()
	fakeClient.Requests = append(fakeClient.Requests, request)
	m.Unlock()
}
