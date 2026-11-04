package services

import "testing"

func TestNormalizeKitchenRole(t *testing.T) {
	for _, role := range []string{"KITCHEN", "kitchen", "CHEF"} {
		if got := normalizeRole(role); got != "CHEF" {
			t.Errorf("normalizeRole(%q) = %q, want CHEF", role, got)
		}
	}
}

func TestLoginRejectsDemoUsersByDefault(t *testing.T) {
	service := NewAuthService("test-secret")

	_, err := service.Login(LoginRequest{
		Username: "admin",
		Password: "LaLuna2025!",
		Role:     "ADMIN",
	})
	if err == nil {
		t.Fatal("expected demo login to be disabled by default")
	}
}
