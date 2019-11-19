package utils

import "time"

type DefaultStoppableTimer struct {
	start time.Time
}

func (timer *DefaultStoppableTimer) Stop() time.Duration {
	return time.Since(timer.start)
}

type DefaultTimer struct {
	start time.Time
}

func (timer *DefaultTimer) Start() StoppableTimer {
	return &DefaultStoppableTimer{
		start: time.Now(),
	}
}

func NewDefaultTimer() *DefaultTimer {
	return &DefaultTimer{}
}
