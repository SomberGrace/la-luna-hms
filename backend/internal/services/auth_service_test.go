package services

import "testing"

func TestNormalizeRole(t *testing.T) {
	tests := []struct {
		role string
		want string
	}{
		{role: "KITCHEN", want: "CHEF"},
		{role: "kitchen", want: "CHEF"},
		{role: " CHEF ", want: "CHEF"},
		{role: "front_desk", want: "FRONT_DESK"},
		{role: "unknown", want: "UNKNOWN"},
	}

	for _, test := range tests {
		if got := normalizeRole(test.role); got != test.want {
			t.Errorf("normalizeRole(%q) = %q, want %q", test.role, got, test.want)
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
