package server

import (
	"github.com/reaandrew/schmokin/service"
	"github.com/reaandrew/schmokin/utils"
)

func MergeResponses(responses []*SchmokinResponse) (result *service.SchmokinResult) {
	result = &service.SchmokinResult{}
	availabilities := []float64{}
	responseTimes := []float64{}
	concurrencyRate := []float64{}
	dateReceiveRates := []float64{}
	dataSendRates := []float64{}
	failedTransactions := []int64{}
	longestTransactions := []int64{}
	shortestTransactions := []int64{}
	successfulTransactions := []int64{}
	totalBytesReceived := []int64{}
	totalBytesSent := []int64{}
	transactions := []int64{}
	transactionRates := []float64{}

	for _, response := range responses {
		availabilities = append(availabilities, response.Availability)
		responseTimes = append(responseTimes, response.AverageResponseTime)
		concurrencyRate = append(concurrencyRate, response.ConcurrencyRate)
		dateReceiveRates = append(dateReceiveRates, response.DataReceiveRate)
		dataSendRates = append(dataSendRates, response.DataSendRate)
		failedTransactions = append(failedTransactions, response.FailedTransactions)
		longestTransactions = append(longestTransactions, response.LongestTransaction)
		shortestTransactions = append(shortestTransactions, response.ShortestTransaction)
		successfulTransactions = append(successfulTransactions, response.SuccessfulTransactions)
		totalBytesReceived = append(totalBytesReceived, int64(response.TotalBytesReceived))
		totalBytesSent = append(totalBytesSent, int64(response.TotalBytesSent))
		transactions = append(transactions, int64(response.Transactions))
		transactionRates = append(transactionRates, response.TransactionRate)
	}

	result.Availability = utils.AverageFloat64(availabilities)
	result.AverageResponseTime = utils.AverageFloat64(responseTimes)
	result.ConcurrencyRate = utils.AverageFloat64(concurrencyRate)
	result.DataReceiveRate = utils.AverageFloat64(dateReceiveRates)
	result.DataSendRate = utils.AverageFloat64(dataSendRates)
	result.FailedTransactions = utils.Sum(failedTransactions)
	result.LongestTransaction = utils.Max(longestTransactions)
	result.ShortestTransaction = utils.Min(shortestTransactions)
	result.SuccessfulTransactions = utils.Sum(successfulTransactions)
	result.TotalBytesReceived = int(utils.Sum(totalBytesReceived))
	result.TotalBytesSent = int(utils.Sum(totalBytesSent))
	result.Transactions = int(utils.Sum(transactions))
	result.TransactionRate = utils.AverageFloat64(transactionRates)
	return result
}
