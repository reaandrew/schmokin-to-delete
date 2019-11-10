package server

import (
	"testing"
	"time"
)

func Test_Something(t *testing.T) {
	go func() {
		StartServer()
	}()

	time.Sleep(1 * time.Second)

	RunClient()
}
