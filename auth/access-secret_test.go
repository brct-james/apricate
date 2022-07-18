package auth_test

import (
	"apricate/auth"
	"testing"
)

func TestGenerateRandomSecureString(t *testing.T) {
	t.Run("empty config file", func(t *testing.T) {
		n := 10
		got, got_error := auth.GenerateRandomSecureString(n)

		if len(got) != n || got_error != nil {
			t.Errorf("n %v, got %q, got_error %q", n, got, got_error)
		}
	})
}