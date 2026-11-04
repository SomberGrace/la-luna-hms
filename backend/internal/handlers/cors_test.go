package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNgrokSkipHeaderAllowedByCORS(t *testing.T) {
	authHandler := NewAuthHandler(nil)
	guestHandler := NewGuestHandler(nil, "")
	tests := []struct {
		name    string
		handler http.Handler
	}{
		{
			name: "shared API handler",
			handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				setCORS(w)
				w.WriteHeader(http.StatusNoContent)
			}),
		},
		{
			name:    "login",
			handler: http.HandlerFunc(authHandler.HandleLogin),
		},
		{
			name:    "registration",
			handler: http.HandlerFunc(authHandler.HandleRegister),
		},
		{
			name:    "guest portal",
			handler: http.HandlerFunc(guestHandler.GetAvailableRooms),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodOptions, "/api/test", nil)
			response := httptest.NewRecorder()
			test.handler.ServeHTTP(response, request)

			if response.Code != http.StatusNoContent {
				t.Fatalf("expected status %d, got %d", http.StatusNoContent, response.Code)
			}
			allowedHeaders := response.Header().Get("Access-Control-Allow-Headers")
			if !strings.Contains(allowedHeaders, "ngrok-skip-browser-warning") {
				t.Errorf("CORS headers %q do not allow ngrok-skip-browser-warning", allowedHeaders)
			}
		})
	}
}

func TestCORSMiddlewareHandlesOperationalPreflight(t *testing.T) {
	called := false
	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodOptions, "/api/v1/restaurant/dockets", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, response.Code)
	}
	if called {
		t.Fatal("preflight request should not reach the application handler")
	}
	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("expected wildcard origin, got %q", got)
	}
	if got := response.Header().Get("Access-Control-Allow-Headers"); !strings.Contains(got, "Authorization") {
		t.Errorf("CORS headers %q do not allow Authorization", got)
	}
}
