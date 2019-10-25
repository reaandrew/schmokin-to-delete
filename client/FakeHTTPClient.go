package client

import "net/http"

type FakeHTTPClient struct {
}

func (fakeClient FakeHTTPClient) Execute(request *http.Request) {

}
