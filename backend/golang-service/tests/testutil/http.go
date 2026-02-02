package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	// Set Gin to test mode to reduce log output
	gin.SetMode(gin.TestMode)
}

// TestRouter creates a new Gin engine for testing.
func TestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())
	return r
}

// TestRequest creates a new HTTP request for testing.
func TestRequest(method, path string, body any) *http.Request {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// TestRequestWithAuth creates a new HTTP request with auth context values.
func TestRequestWithAuth(method, path string, body any, userID uuid.UUID, role string) *http.Request {
	req := TestRequest(method, path, body)
	// Note: For Gin, auth values are typically set via middleware context,
	// not request headers. The handler test should set these in the Gin context.
	return req
}

// PerformRequest executes a request against a Gin router and returns the recorder.
func PerformRequest(router *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := TestRequest(method, path, body)
	router.ServeHTTP(w, req)
	return w
}

// PerformRequestWithContext executes a request with custom context setup.
func PerformRequestWithContext(router *gin.Engine, method, path string, body any, setupCtx func(*gin.Context)) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := TestRequest(method, path, body)

	// Create a context with the setup function
	router.Use(func(c *gin.Context) {
		setupCtx(c)
		c.Next()
	})

	router.ServeHTTP(w, req)
	return w
}

// SetAuthContext is a Gin middleware that sets auth context for testing.
func SetAuthContext(userID uuid.UUID, email, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Set("email", email)
		c.Set("role", role)
		c.Next()
	}
}

// ParseJSONResponse parses the JSON response body into the given struct.
func ParseJSONResponse(t *testing.T, w *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.NewDecoder(w.Body).Decode(v); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}
}

// AssertStatusCode asserts that the response has the expected status code.
func AssertStatusCode(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if w.Code != expected {
		t.Errorf("expected status %d, got %d. Body: %s", expected, w.Code, w.Body.String())
	}
}

// AssertJSONContains asserts that the response body contains the expected key-value pairs.
func AssertJSONContains(t *testing.T, w *httptest.ResponseRecorder, expected map[string]any) {
	t.Helper()

	var actual map[string]any
	if err := json.NewDecoder(w.Body).Decode(&actual); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	for key, expectedVal := range expected {
		actualVal, ok := actual[key]
		if !ok {
			t.Errorf("expected key %q not found in response", key)
			continue
		}
		// Deep comparison for nested structures
		expectedJSON, _ := json.Marshal(expectedVal)
		actualJSON, _ := json.Marshal(actualVal)
		if string(expectedJSON) != string(actualJSON) {
			t.Errorf("key %q: expected %v, got %v", key, expectedVal, actualVal)
		}
	}
}
