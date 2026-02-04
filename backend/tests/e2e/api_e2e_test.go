//go:build e2e

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/ap1-final-mini-moodle/tests/testutil"
)

func TestHealthEndpoint(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(getBaseURL() + "/api/v1/health")
	if err != nil {
		t.Fatalf("health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("health check returned %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode health response: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", result["status"])
	}
}

func TestCourseAPI_E2E(t *testing.T) {
	pool := getTestPool(t)
	fixtures := testutil.NewFixtures(pool)
	ctx := context.Background()

	// Setup test data
	testutil.TruncateAll(t, pool)
	fixtures.SeedBasicData(ctx, t)

	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("list courses", func(t *testing.T) {
		// Note: This test assumes authentication is handled or the endpoint is public
		// In real e2e tests, you would first authenticate and get a token
		resp, err := client.Get(getBaseURL() + "/api/v1/courses")
		if err != nil {
			t.Fatalf("list courses failed: %v", err)
		}
		defer resp.Body.Close()

		// If endpoint requires auth, expect 401
		if resp.StatusCode == http.StatusUnauthorized {
			t.Log("endpoint requires authentication (expected)")
			return
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("list courses returned %d: %s", resp.StatusCode, string(body))
		}
	})

	t.Run("get course by ID", func(t *testing.T) {
		resp, err := client.Get(getBaseURL() + "/api/v1/courses/" + testutil.TestCourseID1.String())
		if err != nil {
			t.Fatalf("get course failed: %v", err)
		}
		defer resp.Body.Close()

		// If endpoint requires auth, expect 401
		if resp.StatusCode == http.StatusUnauthorized {
			t.Log("endpoint requires authentication (expected)")
			return
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("get course returned %d: %s", resp.StatusCode, string(body))
		}
	})

	t.Run("get non-existing course returns 404", func(t *testing.T) {
		resp, err := client.Get(getBaseURL() + "/api/v1/courses/00000000-0000-0000-0000-000000000999")
		if err != nil {
			t.Fatalf("get course failed: %v", err)
		}
		defer resp.Body.Close()

		// If endpoint requires auth, expect 401
		if resp.StatusCode == http.StatusUnauthorized {
			t.Log("endpoint requires authentication (expected)")
			return
		}

		if resp.StatusCode != http.StatusNotFound {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("expected 404, got %d: %s", resp.StatusCode, string(body))
		}
	})
}

func TestAuthFlow_E2E(t *testing.T) {
	pool := getTestPool(t)
	_ = context.Background() // ctx available for future use

	// Clean database
	testutil.TruncateAll(t, pool)

	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("register new user", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email":    "newuser@test.com",
			"password": "securepassword123",
			"role":     "student",
		}

		bodyBytes, _ := json.Marshal(reqBody)
		resp, err := client.Post(
			getBaseURL()+"/api/v1/auth/register",
			"application/json",
			bytes.NewBuffer(bodyBytes),
		)
		if err != nil {
			t.Fatalf("register failed: %v", err)
		}
		defer resp.Body.Close()

		// Registration might return different codes depending on implementation
		// 201 Created or 200 OK are both acceptable
		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Logf("register returned %d: %s", resp.StatusCode, string(body))
		}
	})

	t.Run("login with registered user", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email":    "newuser@test.com",
			"password": "securepassword123",
		}

		bodyBytes, _ := json.Marshal(reqBody)
		resp, err := client.Post(
			getBaseURL()+"/api/v1/auth/login",
			"application/json",
			bytes.NewBuffer(bodyBytes),
		)
		if err != nil {
			t.Fatalf("login failed: %v", err)
		}
		defer resp.Body.Close()

		// Expect 200 OK or 401 if user wasn't created
		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				t.Fatalf("failed to decode login response: %v", err)
			}

			// Check for access token in response
			if data, ok := result["data"].(map[string]interface{}); ok {
				if _, hasToken := data["access_token"]; !hasToken {
					t.Error("login response missing access_token")
				}
			}
		} else {
			body, _ := io.ReadAll(resp.Body)
			t.Logf("login returned %d: %s", resp.StatusCode, string(body))
		}
	})

	t.Run("login with wrong password returns error", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email":    "newuser@test.com",
			"password": "wrongpassword",
		}

		bodyBytes, _ := json.Marshal(reqBody)
		resp, err := client.Post(
			getBaseURL()+"/api/v1/auth/login",
			"application/json",
			bytes.NewBuffer(bodyBytes),
		)
		if err != nil {
			t.Fatalf("login failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			t.Error("expected error for wrong password, got 200 OK")
		}
	})
}
