package utils

import "time"

type DefaultTimer struct {
	start time.Time
}

func (timer *DefaultTimer) Start() {
	timer.start = time.Now()
}

func (timer *DefaultTimer) Stop() time.Duration {
	return time.Since(timer.start)
}
