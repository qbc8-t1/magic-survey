package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloAnswerAPI(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/api/v1/answers/hello")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode, "expected HTTP 200 OK")

	// Validate response body
	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, "Hello from Answer Service!", body["message"])
}

func TestCreateAnswerAPI(t *testing.T) {
	answer := map[string]interface{}{
		"submission_id": 1,
		"question_id":   1,
		"option_id":     2,
		"text":          "Sample answer text",
	}
	body, err := json.Marshal(answer)
	assert.NoError(t, err)

	resp, err := http.Post("http://localhost:8080/api/v1/answers", "application/json", bytes.NewReader(body))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "expected HTTP 201 Created")

	// Validate response body
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, answer["submission_id"], response["submission_id"])
	assert.Equal(t, answer["question_id"], response["question_id"])
}

func TestGetAnswerAPI(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/api/v1/answers/1")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode, "expected HTTP 200 OK")

	// Validate response body
	var answer map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&answer)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), answer["id"]) // JSON unmarshals numbers as float64
}

func TestUpdateAnswerAPI(t *testing.T) {
	updateData := map[string]interface{}{
		"text": "Updated answer text",
	}
	body, err := json.Marshal(updateData)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/answers/1", bytes.NewReader(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode, "expected HTTP 200 OK")

	// Validate response body
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, updateData["text"], response["text"])
}

func TestDeleteAnswerAPI(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/answers/1", nil)
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusNoContent, resp.StatusCode, "expected HTTP 204 No Content")
}
