package utils

import "time"

type FakeTimer struct {
	elapsed time.Duration
}

func (timer *FakeTimer) SetElapsed(duration time.Duration) {
	timer.elapsed = duration
}

func (timer *FakeTimer) Start() {

}

func (timer *FakeTimer) Stop() time.Duration {
	return timer.elapsed
}

func NewFakeTimer(duration time.Duration) *FakeTimer {
	timer := &FakeTimer{}
	timer.SetElapsed(duration)
	return timer
}
