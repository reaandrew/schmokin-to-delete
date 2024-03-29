package service_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	schmokinHTTP "github.com/reaandrew/schmokin/infrastructure/http"
	"github.com/reaandrew/schmokin/service"
	"github.com/reaandrew/schmokin/utils"
	"github.com/stretchr/testify/assert"
)

type SchmokinServiceTransactionTestCase struct {
	Urls                 int
	Workers              int
	Iterations           int
	ExpectedTransactions int
}

// All these tests need to be on the SchmokinService not the SchmokinService
// The Schmokin CLI should be tested to ensure it invokes the correct proxy.

func Test_SchmokinServiceReturnNumberOfTransactions(t *testing.T) {
	cases := []SchmokinServiceTransactionTestCase{
		{Urls: 1, Workers: 1, Iterations: 1, ExpectedTransactions: 1},
		{Urls: 2, Workers: 1, Iterations: 1, ExpectedTransactions: 2},
		{Urls: 1, Workers: 2, Iterations: 1, ExpectedTransactions: 2},
		{Urls: 2, Workers: 2, Iterations: 1, ExpectedTransactions: 4},
		{Urls: 1, Workers: 2, Iterations: 2, ExpectedTransactions: 4},
		{Urls: 3, Workers: 2, Iterations: 2, ExpectedTransactions: 4},
		{Urls: 2, Workers: 2, Iterations: 3, ExpectedTransactions: 6},
		{Urls: 5, Workers: 100, Iterations: 5, ExpectedTransactions: 500},
	}

	httpClient := schmokinHTTP.NewFakeClient()
	for _, currentTestCase := range cases {
		testCase := currentTestCase
		t.Run(fmt.Sprintf("Test_SchmokinServiceReturnNumberOfTransactions_Urls_%v_Workers_%v_Iterations_%v_Returns_%v_Transactions",
			testCase.Urls,
			testCase.Workers,
			testCase.Iterations,
			testCase.ExpectedTransactions), func(t *testing.T) {
			lines := utils.CreateRandomLines(testCase.Urls)
			schmokinService := service.NewSchmokinServiceBuilder().
				SetClient(httpClient).
				SetWorkers(testCase.Workers).
				SetIterations(testCase.Iterations).
				Build()
			result := schmokinService.Execute(lines)

			assert.Equal(t, testCase.ExpectedTransactions, result.Transactions)
		})
	}
}

type SchmokinServiceAvailabilityTestCase struct {
	StatusCodes          []int
	ExpectedAvailability float64
}

func Test_SchmokinServiceReturnsAvailability(t *testing.T) {
	cases := []SchmokinServiceAvailabilityTestCase{
		{StatusCodes: []int{200, 200, 500, 500}, ExpectedAvailability: float64(0.5)},
		{StatusCodes: []int{200, 200}, ExpectedAvailability: float64(1)},
		{StatusCodes: []int{200, 201, 202}, ExpectedAvailability: float64(1)},
		{StatusCodes: []int{200, 200, 404, 500}, ExpectedAvailability: float64(0.5)},
		{StatusCodes: []int{500, 500, 500, 500}, ExpectedAvailability: float64(0)},
	}

	for _, currentTestCase := range cases {
		testCase := currentTestCase
		t.Run(fmt.Sprintf("Test_SchmokinServiceReturnAvailabilityOf%v%%", testCase.ExpectedAvailability*100), func(t *testing.T) {
			lines := utils.CreateRandomLines(len(testCase.StatusCodes))
			httpClient := schmokinHTTP.NewFakeClient()
			schmokinService := service.NewSchmokinServiceBuilder().
				SetClient(httpClient).
				Build()
			count := 0
			httpClient.Interceptor = func(response *http.Response) {
				response.StatusCode = testCase.StatusCodes[count]
				count++
			}
			result := schmokinService.Execute(lines)

			assert.Equal(t, testCase.ExpectedAvailability, result.Availability)
		})
	}
}

func Test_SchmokinServiceReturnsElapsedTime(t *testing.T) {
	lines := utils.CreateRandomLines(1)
	expectedElapsed := 100 * time.Second
	timer := &utils.FakeTimer{}
	timer.SetElapsed(expectedElapsed)
	httpClient := schmokinHTTP.NewFakeClient()
	schmokinService := service.NewSchmokinServiceBuilder().
		SetTimer(timer).
		SetClient(httpClient).
		Build()
	result := schmokinService.Execute(lines)

	assert.Equal(t, expectedElapsed, result.ElapsedTime)
}

func Test_SchmokinServiceReturnsTotalBytesSent(t *testing.T) {
	lines := utils.CreateRandomLines(1)
	httpClient := schmokinHTTP.NewFakeClient()
	schmokinService := service.NewSchmokinServiceBuilder().
		SetClient(httpClient).
		Build()
	result := schmokinService.Execute(lines)

	// This is the size of one request dumped
	assert.Equal(t, result.TotalBytesSent, 96)
}

func Test_SchmokinServiceReturnsTotalBytesReceived(t *testing.T) {
	lines := utils.CreateRandomLines(1)
	httpClient := schmokinHTTP.NewFakeClient()
	schmokinService := service.NewSchmokinServiceBuilder().
		SetClient(httpClient).
		Build()
	result := schmokinService.Execute(lines)

	// This is the size of one request dumped
	assert.Equal(t, result.TotalBytesReceived, 38)
}

func Test_SchmokinServiceReturnsAverageResponseTime(t *testing.T) {
	lines := utils.CreateRandomLines(1)
	expectedDuration := 1 * time.Minute
	timer := utils.NewFakeTimer(expectedDuration)
	httpClient := schmokinHTTP.NewFakeClient()
	schmokinService := service.NewSchmokinServiceBuilder().
		SetTimer(timer).
		SetClient(httpClient).
		Build()
	result := schmokinService.Execute(lines)

	// This is the size of one request dumped
	assert.Equal(t, float64(expectedDuration), result.AverageResponseTime)
}
