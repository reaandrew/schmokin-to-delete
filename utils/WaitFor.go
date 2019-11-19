package utils

import (
	"time"
)

type WaitUtil struct {
	Timeout time.Duration
	Backoff time.Duration
}

func (waitUtil WaitUtil) Wait(waiter func() bool) {
	start := time.Now()
	for {
		if time.Since(start) > waitUtil.Timeout {
			panic("Timed out")
		}
		if waiter() {
			return
		}
		time.Sleep(waitUtil.Backoff)
	}
}
