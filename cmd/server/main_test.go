package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestServerRoutes(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test server
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test the config endpoint
	t.Run("GET /api/config", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/config", nil)
		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check that the response contains expected fields
		_, exists := response["defaultModel"]
		assert.True(t, exists)

		_, exists = response["hasCompletedOnboarding"]
		assert.True(t, exists)
	})

	// Test the files endpoint
	t.Run("GET /api/files", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/files?path=.", nil)
		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check that the response contains entries
		_, exists := response["entries"]
		assert.True(t, exists)
	})

	// Test the chat endpoint
	t.Run("POST /api/chat", func(t *testing.T) {
		// Skip this test as it requires an API key
		t.Skip("Skipping chat test as it requires an API key")

		w := httptest.NewRecorder()
		reqBody := `{"sessionId":"test-session","message":"Hello"}`
		req, _ := http.NewRequest("POST", "/api/chat", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check that the response contains expected fields
		_, exists := response["response"]
		assert.True(t, exists)
	})
}
