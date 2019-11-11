package client_test

import (
	"net/http"
	"testing"

	"github.com/reaandrew/surge/client"
	"github.com/stretchr/testify/assert"
)

func Test_WhenNoServerExists(t *testing.T) {
	request, _ := http.NewRequest("GET", "http://localhost:45000", nil)
	httpClient := client.NewDefaultHttpClient()
	response, err := httpClient.Execute(request)
	assert.Nil(t, response)
	assert.NotNil(t, err)
}
