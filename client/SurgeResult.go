package client

import "time"

type Result struct {
	Transactions           int
	Availability           float64
	ElapsedTime            time.Duration
	AverageResponseTime    float64
	TotalBytesSent         int
	TotalBytesReceived     int
	TransactionRate        float64
	ConcurrencyRate        float64
	DataSendRate           float64
	DataReceiveRate        float64
	SuccessfulTransactions int64
	FailedTransactions     int64
	LongestTransaction     int64
	ShortestTransaction    int64
}
