package client

import "time"

type Result struct {
	Transactions   int
	Availability   float64
	ElapsedTime    time.Duration
	TotalBytesSent int
}
