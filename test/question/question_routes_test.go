package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloQuestionAPI(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/api/v1/questions/hello")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode, "expected HTTP 200 OK")

	// Validate response body
	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, "Hello from Question Service!", body["message"])
}

func TestCreateQuestionAPI(t *testing.T) {
	question := map[string]interface{}{
		"title":            "Sample Question",
		"type":             "multioption",
		"questionnaire_id": 1,
		"order":            1,
		"options": []map[string]string{
			{"text": "Option 1"},
			{"text": "Option 2"},
		},
	}
	body, err := json.Marshal(question)
	assert.NoError(t, err)

	resp, err := http.Post("http://localhost:8080/api/v1/questions", "application/json", bytes.NewReader(body))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "expected HTTP 201 Created")

	// Validate response body
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, question["title"], response["title"])
	assert.Equal(t, question["type"], response["type"])
}

func TestGetQuestionAPI(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/api/v1/questions/1")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode, "expected HTTP 200 OK")

	// Validate response body
	var question map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&question)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), question["id"]) // JSON unmarshals numbers as float64
}

func TestUpdateQuestionAPI(t *testing.T) {
	updateData := map[string]interface{}{
		"title": "Updated Question Title",
	}
	body, err := json.Marshal(updateData)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/questions/1", bytes.NewReader(body))
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
	assert.Equal(t, updateData["title"], response["title"])
}

func TestDeleteQuestionAPI(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/questions/1", nil)
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusNoContent, resp.StatusCode, "expected HTTP 204 No Content")
}

func TestGetAllQuestionsAPI(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/api/v1/questions")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode, "expected HTTP 200 OK")

	// Validate response body
	var questions []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&questions)
	assert.NoError(t, err)
	assert.Greater(t, len(questions), 0, "expected at least one question")
}
