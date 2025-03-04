package forecaster

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Mock HTTP Client
type mockHTTPClient struct {
	handler http.HandlerFunc
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	m.handler(rec, req)
	return rec.Result(), nil
}

// Helper function to create a test WeatherClient
func newTestClient(handler http.HandlerFunc) *WeatherClient {
	return &WeatherClient{
		client:  &mockHTTPClient{handler: handler},
		baseURL: "http://mockapi.com",
		key:     "test-api-key",
	}
}

// Test Successful API Call
func TestGetWeather_Success(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"temperature": 25, "condition": "Sunny", "weather": [{"main": "Clear"}]}`)) // Simulated JSON response
	}

	client := newTestClient(mockHandler)
	weather, errResp := client.GetWeather(51.5074, -0.1278,nil) // London Coordinates

	assert.Nil(t, errResp, "Error response should be nil")
	assert.NotNil(t, weather, "Weather response should not be nil")
}

// Test Timeout Handling
func TestGetWeather_Timeout(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"temperature": 22, "condition": "Cloudy"}`))
	}

	client := newTestClient(mockHandler)

	start := time.Now()
	weather, errResp := client.GetWeather(51.5074, -0.1278,nil) // London

	duration := time.Since(start)

	assert.Nil(t, weather, "Weather response should be nil on timeout")
	assert.NotNil(t, errResp, "Error response should not be nil")
	assert.Equal(t, "cannot do request, deadline expired", errResp.Message)
	assert.Less(t, duration.Milliseconds(), int64(250), "Request should timeout before 2s")
}

// Test API Returning Non-200 Response
func TestGetWeather_ServerError(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Cod": "500", "Message": "Internal Server Error"}`))
	}

	client := newTestClient(mockHandler)
	weather, errResp := client.GetWeather(51.5074, -0.1278,nil)

	assert.Nil(t, weather, "Weather response should be nil")
	assert.NotNil(t, errResp, "Error response should not be nil")
	assert.Equal(t, "Internal Server Error", errResp.Message)
}

// Test HTTP Request Failure
func TestGetWeather_RequestFailure(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"Cod": "502", "Message": "Bad Gateway"}`))
	}

	client := newTestClient(mockHandler)
	weather, errResp := client.GetWeather(51.5074, -0.1278,nil)

	assert.Nil(t, weather, "Weather response should be nil")
	assert.NotNil(t, errResp, "Error response should not be nil")
	assert.Equal(t, "Bad Gateway", errResp.Message)
}

// Test Invalid JSON Response
func TestGetWeather_InvalidJSON(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid json}`)) // Malformed JSON
	}

	client := newTestClient(mockHandler)
	weather, errResp := client.GetWeather(51.5074, -0.1278,nil)

	assert.Nil(t, weather, "Weather response should be nil")
	assert.NotNil(t, errResp, "Error response should not be nil")
	assert.Equal(t, "cannot decode response", errResp.Message)
}

