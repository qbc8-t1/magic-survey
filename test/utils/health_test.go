package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthAPI(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/health")
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	defer resp.Body.Close()

	// check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode, "expected HTTP 200 OK")
}
