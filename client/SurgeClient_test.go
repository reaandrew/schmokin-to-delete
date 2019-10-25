package client_test

import (
	"fmt"
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
				HttpClient:  &client.FakeHTTPClient{},
				Iterations:  testCase.Iterations,
			}
			transactions, err := client.Run()

			assert.Nil(t, err)
			assert.Equal(t, testCase.ExpectedTransactions, transactions)
		})
	}
}
