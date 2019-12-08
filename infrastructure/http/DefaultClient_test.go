package http_test

import (
	"net/http"
	"testing"

	schmokinHTTP "github.com/reaandrew/schmokin/infrastructure/http"

	"github.com/stretchr/testify/assert"
)

func Test_WhenNoServerExists(t *testing.T) {
	request, _ := http.NewRequest("GET", "http://localhost:45000", nil)
	httpClient := schmokinHTTP.NewDefaultClient()
	response, err := httpClient.Execute(request)

	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	assert.Nil(t, response)
	assert.NotNil(t, err)
}
