package comparison_test

import (
	"testing"

	"apricate/comparison"
)

func TestStringSlicesEqual(t *testing.T) {
	t.Run("return false when different sizes", func(t *testing.T) {
		a := []string{"is", "same"}
		b := []string{"is", "not", "same"}
		got := comparison.StringSlicesEqual(a, b)
		want := false

		if got != want {
			t.Errorf("a %q, b %q, got %v, want %v", a, b, got, want)
		}
	})
	t.Run("return false when different elements", func(t *testing.T) {
		a := []string{"is", "same"}
		b := []string{"not", "same"}
		got := comparison.StringSlicesEqual(a, b)
		want := false

		if got != want {
			t.Errorf("a %q, b %q, got %v, want %v", a, b, got, want)
		}
	})
	t.Run("return true when same size and elements", func(t *testing.T) {
		a := []string{"is", "same"}
		b := []string{"is", "same"}
		got := comparison.StringSlicesEqual(a, b)
		want := true

		if got != want {
			t.Errorf("a %q, b %q, got %v, want %v", a, b, got, want)
		}
	})
}