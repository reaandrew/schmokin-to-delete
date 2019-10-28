package client_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/reaandrew/surge/client"
	"github.com/reaandrew/surge/utils"
	"github.com/stretchr/testify/assert"
)

type SurgeClientTransactionTestCase struct {
	Urls                 int
	Workers              int
	Iterations           int
	ExpectedTransactions int
}

func Test_SurgeClientReturnNumberOfTransactions(t *testing.T) {
	cases := []SurgeClientTransactionTestCase{
		SurgeClientTransactionTestCase{Urls: 1, Workers: 1, Iterations: 1, ExpectedTransactions: 1},
		SurgeClientTransactionTestCase{Urls: 2, Workers: 1, Iterations: 1, ExpectedTransactions: 2},
		SurgeClientTransactionTestCase{Urls: 1, Workers: 2, Iterations: 1, ExpectedTransactions: 2},
		SurgeClientTransactionTestCase{Urls: 2, Workers: 2, Iterations: 1, ExpectedTransactions: 4},
		SurgeClientTransactionTestCase{Urls: 1, Workers: 2, Iterations: 2, ExpectedTransactions: 4},
		SurgeClientTransactionTestCase{Urls: 3, Workers: 2, Iterations: 2, ExpectedTransactions: 4},
		SurgeClientTransactionTestCase{Urls: 2, Workers: 2, Iterations: 3, ExpectedTransactions: 6},
		SurgeClientTransactionTestCase{Urls: 5, Workers: 100, Iterations: 5, ExpectedTransactions: 500},
	}

	for _, testCase := range cases {
		t.Run(fmt.Sprintf("Test_SurgeClientReturnNumberOfTransactions_Urls_%v_Workers_%v_Iterations_%v_Returns_%v_Transactions",
			testCase.Urls,
			testCase.Workers,
			testCase.Iterations,
			testCase.ExpectedTransactions), func(t *testing.T) {
			file := utils.CreateRandomHttpTestFile(testCase.Urls)
			surgeClient := client.NewSurgeClientBuilder().
				SetURLFilePath(file.Name()).
				SetWorkers(testCase.Workers).
				SetHTTPClient(client.NewFakeHTTPClient()).
				SetIterations(testCase.Iterations).
				Build()
			result, err := surgeClient.Run()

			assert.Nil(t, err)
			assert.Equal(t, testCase.ExpectedTransactions, result.Transactions)
		})
	}
}

type SurgeClientAvailabilityTestCase struct {
	StatusCodes          []int
	ExpectedAvailability float64
}

func Test_SurgeClientReturnsAvailability(t *testing.T) {
	cases := []SurgeClientAvailabilityTestCase{
		SurgeClientAvailabilityTestCase{StatusCodes: []int{200, 200, 500, 500}, ExpectedAvailability: float64(0.5)},
		SurgeClientAvailabilityTestCase{StatusCodes: []int{200, 200}, ExpectedAvailability: float64(1)},
		SurgeClientAvailabilityTestCase{StatusCodes: []int{200, 201, 202}, ExpectedAvailability: float64(1)},
		SurgeClientAvailabilityTestCase{StatusCodes: []int{200, 200, 404, 500}, ExpectedAvailability: float64(0.5)},
		SurgeClientAvailabilityTestCase{StatusCodes: []int{500, 500, 500, 500}, ExpectedAvailability: float64(0)},
	}

	for _, testCase := range cases {
		t.Run(fmt.Sprintf("Test_SurgeClientReturnAvailabilityOf%v%%", testCase.ExpectedAvailability*100), func(t *testing.T) {
			file := utils.CreateRandomHttpTestFile(len(testCase.StatusCodes))
			httpClient := client.NewFakeHTTPClient()
			surgeClient := client.NewSurgeClientBuilder().
				SetURLFilePath(file.Name()).
				SetHTTPClient(httpClient).
				Build()
			count := 0
			httpClient.Interceptor = func(response *http.Response) {
				response.StatusCode = testCase.StatusCodes[count]
				count++
			}
			result, err := surgeClient.Run()

			assert.Nil(t, err)
			assert.Equal(t, testCase.ExpectedAvailability, result.Availability)
		})
	}
}

func Test_SurgeClientReturnsElapsedTime(t *testing.T) {

}
