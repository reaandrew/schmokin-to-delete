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
			client := client.Surge{
				UrlFilePath: file.Name(),
				WorkerCount: testCase.Workers,
				HttpClient:  client.NewFakeHTTPClient(),
				Iterations:  testCase.Iterations,
			}
			result, err := client.Run()

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
	}

	for _, testCase := range cases {
		t.Run(fmt.Sprintf("Test_SurgeClientReturnAvailabilityOf%v%%", testCase.ExpectedAvailability*100), func(t *testing.T) {
			file := utils.CreateRandomHttpTestFile(len(testCase.StatusCodes))
			httpClient := client.NewFakeHTTPClient()
			client := client.Surge{
				UrlFilePath: file.Name(),
				WorkerCount: 1,
				HttpClient:  httpClient,
				Iterations:  1,
			}
			count := 0
			httpClient.Interceptor = func(response *http.Response) {
				response.StatusCode = testCase.StatusCodes[count]
				count++
			}
			result, err := client.Run()

			assert.Nil(t, err)
			assert.Equal(t, testCase.ExpectedAvailability, result.Availability)
		})
	}
}

func Test_SurgeClientReturnAvailabilityOf1(t *testing.T) {
	file := utils.CreateRandomHttpTestFile(1)
	client := client.Surge{
		UrlFilePath: file.Name(),
		WorkerCount: 1,
		HttpClient:  client.NewFakeHTTPClient(),
		Iterations:  1,
	}
	result, err := client.Run()

	assert.Nil(t, err)
	assert.Equal(t, result.Availability, float64(1))
}

func Test_SurgeClientReturnsAvailabilityOf0_5(t *testing.T) {
	file := utils.CreateRandomHttpTestFile(10)
	httpClient := client.NewFakeHTTPClient()
	client := client.Surge{
		UrlFilePath: file.Name(),
		WorkerCount: 1,
		HttpClient:  httpClient,
		Iterations:  1,
	}
	count := 0
	httpClient.Interceptor = func(response *http.Response) {
		if count%2 == 0 {
			response.StatusCode = 500
		}
		count++
	}
	result, err := client.Run()

	assert.Nil(t, err)
	assert.Equal(t, result.Availability, float64(0.5))
}
