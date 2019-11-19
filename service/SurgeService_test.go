package service_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	surgeHTTP "github.com/reaandrew/surge/infrastructure/http"
	"github.com/reaandrew/surge/service"
	"github.com/reaandrew/surge/utils"
	"github.com/stretchr/testify/assert"
)

type SurgeServiceTransactionTestCase struct {
	Urls                 int
	Workers              int
	Iterations           int
	ExpectedTransactions int
}

// All these tests need to be on the SurgeService not the SurgeService
// The Surge CLI should be tested to ensure it invokes the correct proxy.

func Test_SurgeServiceReturnNumberOfTransactions(t *testing.T) {
	cases := []SurgeServiceTransactionTestCase{
		{Urls: 1, Workers: 1, Iterations: 1, ExpectedTransactions: 1},
		{Urls: 2, Workers: 1, Iterations: 1, ExpectedTransactions: 2},
		{Urls: 1, Workers: 2, Iterations: 1, ExpectedTransactions: 2},
		{Urls: 2, Workers: 2, Iterations: 1, ExpectedTransactions: 4},
		{Urls: 1, Workers: 2, Iterations: 2, ExpectedTransactions: 4},
		{Urls: 3, Workers: 2, Iterations: 2, ExpectedTransactions: 4},
		{Urls: 2, Workers: 2, Iterations: 3, ExpectedTransactions: 6},
		{Urls: 5, Workers: 100, Iterations: 5, ExpectedTransactions: 500},
	}

	httpClient := surgeHTTP.NewFakeClient()
	for _, currentTestCase := range cases {
		testCase := currentTestCase
		t.Run(fmt.Sprintf("Test_SurgeServiceReturnNumberOfTransactions_Urls_%v_Workers_%v_Iterations_%v_Returns_%v_Transactions",
			testCase.Urls,
			testCase.Workers,
			testCase.Iterations,
			testCase.ExpectedTransactions), func(t *testing.T) {
			lines := utils.CreateRandomLines(testCase.Urls)
			surgeService := service.NewSurgeServiceBuilder().
				SetClient(httpClient).
				SetWorkers(testCase.Workers).
				SetIterations(testCase.Iterations).
				Build()
			result := surgeService.Execute(lines)

			assert.Equal(t, testCase.ExpectedTransactions, result.Transactions)
		})
	}
}

type SurgeServiceAvailabilityTestCase struct {
	StatusCodes          []int
	ExpectedAvailability float64
}

func Test_SurgeServiceReturnsAvailability(t *testing.T) {
	cases := []SurgeServiceAvailabilityTestCase{
		{StatusCodes: []int{200, 200, 500, 500}, ExpectedAvailability: float64(0.5)},
		{StatusCodes: []int{200, 200}, ExpectedAvailability: float64(1)},
		{StatusCodes: []int{200, 201, 202}, ExpectedAvailability: float64(1)},
		{StatusCodes: []int{200, 200, 404, 500}, ExpectedAvailability: float64(0.5)},
		{StatusCodes: []int{500, 500, 500, 500}, ExpectedAvailability: float64(0)},
	}

	for _, currentTestCase := range cases {
		testCase := currentTestCase
		t.Run(fmt.Sprintf("Test_SurgeServiceReturnAvailabilityOf%v%%", testCase.ExpectedAvailability*100), func(t *testing.T) {
			lines := utils.CreateRandomLines(len(testCase.StatusCodes))
			httpClient := surgeHTTP.NewFakeClient()
			surgeService := service.NewSurgeServiceBuilder().
				SetClient(httpClient).
				Build()
			count := 0
			httpClient.Interceptor = func(response *http.Response) {
				response.StatusCode = testCase.StatusCodes[count]
				count++
			}
			result := surgeService.Execute(lines)

			assert.Equal(t, testCase.ExpectedAvailability, result.Availability)
		})
	}
}

func Test_SurgeServiceReturnsElapsedTime(t *testing.T) {
	lines := utils.CreateRandomLines(1)
	expectedElapsed := 100 * time.Second
	timer := &utils.FakeTimer{}
	timer.SetElapsed(expectedElapsed)
	httpClient := surgeHTTP.NewFakeClient()
	surgeService := service.NewSurgeServiceBuilder().
		SetTimer(timer).
		SetClient(httpClient).
		Build()
	result := surgeService.Execute(lines)

	assert.Equal(t, expectedElapsed, result.ElapsedTime)
}

func Test_SurgeServiceReturnsTotalBytesSent(t *testing.T) {
	lines := utils.CreateRandomLines(1)
	httpClient := surgeHTTP.NewFakeClient()
	surgeService := service.NewSurgeServiceBuilder().
		SetClient(httpClient).
		Build()
	result := surgeService.Execute(lines)

	// This is the size of one request dumped
	assert.Equal(t, result.TotalBytesSent, 96)
}

func Test_SurgeServiceReturnsTotalBytesReceived(t *testing.T) {
	lines := utils.CreateRandomLines(1)
	httpClient := surgeHTTP.NewFakeClient()
	surgeService := service.NewSurgeServiceBuilder().
		SetClient(httpClient).
		Build()
	result := surgeService.Execute(lines)

	// This is the size of one request dumped
	assert.Equal(t, result.TotalBytesReceived, 38)
}

func Test_SurgeServiceReturnsAverageResponseTime(t *testing.T) {
	lines := utils.CreateRandomLines(1)
	expectedDuration := 1 * time.Minute
	timer := utils.NewFakeTimer(expectedDuration)
	httpClient := surgeHTTP.NewFakeClient()
	surgeService := service.NewSurgeServiceBuilder().
		SetTimer(timer).
		SetClient(httpClient).
		Build()
	result := surgeService.Execute(lines)

	// This is the size of one request dumped
	assert.Equal(t, float64(expectedDuration), result.AverageResponseTime)
}
