package utils

import (
	"time"
)

type WaitUtil struct {
	Timeout time.Duration
	Backoff time.Duration
}

func (self WaitUtil) Wait(waiter func() bool) {
	start := time.Now()
	for {
		if time.Since(start) > self.Timeout {
			panic("Timed out")
		}
		if waiter() {
			return
		}
		time.Sleep(self.Backoff)
	}
}
