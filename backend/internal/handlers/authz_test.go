package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestRequireRoleAllowsMatchingRole(t *testing.T) {
	secret := "test-secret"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  "user-42",
		"username": "frontdesk",
		"role":     "FRONT_DESK",
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	jwtString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	called := false
	h := RequireRole(secret, "FRONT_DESK", "MANAGER")(func(w http.ResponseWriter, r *http.Request) {
		called = true
		user := UserFromRequest(r)
		if user == nil || user.Role != "FRONT_DESK" {
			t.Fatalf("expected FRONT_DESK user, got %+v", user)
		}
	})

	req := httptest.NewRequest("GET", "/api/v1/frontdesk/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+jwtString)
	res := httptest.NewRecorder()
	http.HandlerFunc(h).ServeHTTP(res, req)

	if !called {
		t.Fatal("expected handler to be called")
	}
	if res.Code != 200 {
		t.Fatalf("expected 200 but got %d", res.Code)
	}
}

func TestRequireRoleRejectsForbiddenRole(t *testing.T) {
	secret := "test-secret"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  "guest-1",
		"username": "guest201",
		"role":     "GUEST",
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	jwtString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	called := false
	h := RequireRole(secret, "FRONT_DESK", "MANAGER")(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	req := httptest.NewRequest("GET", "/api/v1/frontdesk/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+jwtString)
	res := httptest.NewRecorder()
	http.HandlerFunc(h).ServeHTTP(res, req)

	if called {
		t.Fatal("expected handler not to run")
	}
	if res.Code != http.StatusForbidden {
		t.Fatalf("expected 403 but got %d", res.Code)
	}
}
