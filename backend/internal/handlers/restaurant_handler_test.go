package handlers

import "testing"

func TestIsUUID(t *testing.T) {
	tests := []struct {
		value string
		valid bool
	}{
		{value: "acde4806-c3fc-4ed0-910b-b0c6eb62c85d", valid: true},
		{value: "ACDE4806-C3FC-4ED0-910B-B0C6EB62C85D", valid: true},
		{value: "acde4806c3fc4ed0910bb0c6eb62c85d", valid: false},
		{value: "acde4806-c3fc-4ed0-910b-b0c6eb62c85x", valid: false},
	}

	for _, test := range tests {
		if got := isUUID(test.value); got != test.valid {
			t.Errorf("isUUID(%q) = %t, want %t", test.value, got, test.valid)
		}
	}
}
