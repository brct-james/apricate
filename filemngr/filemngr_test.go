package filemngr_test

import (
	"apricate/filemngr"
	"testing"
)

func TestGetKeyFromLines(t *testing.T) {
	t.Run("return empty string when not found", func(t *testing.T) {
		search_key := "key"
		lines := []string{"not key value", "incorrect_key=value"}
		got := filemngr.GetKeyFromLines(search_key, lines)
		want := ""

		if got != want {
			t.Errorf("search_key %q, lines %q, got %q, want %q", search_key, lines, got, want)
		}
	})
	t.Run("return empty string when lines empty", func(t *testing.T) {
		search_key := "key"
		lines := []string{}
		got := filemngr.GetKeyFromLines(search_key, lines)
		want := ""

		if got != want {
			t.Errorf("search_key %q, lines %q, got %q, want %q", search_key, lines, got, want)
		}
	})
	t.Run("return value string when found", func(t *testing.T) {
		search_key := "key"
		lines := []string{"not key value", "key=value"}
		got := filemngr.GetKeyFromLines(search_key, lines)
		want := "value"

		if got != want {
			t.Errorf("search_key %q, lines %q, got %q, want %q", search_key, lines, got, want)
		}
	})
	t.Run("return value string when found and value contains =", func(t *testing.T) {
		search_key := "key"
		lines := []string{"key=value=test"}
		got := filemngr.GetKeyFromLines(search_key, lines)
		want := "value=test"

		if got != want {
			t.Errorf("search_key %q, lines %q, got %q, want %q", search_key, lines, got, want)
		}
	})
	t.Run("return empty string for notkey=key=value", func(t *testing.T) {
		search_key := "key"
		lines := []string{"notkey=key=value"}
		got := filemngr.GetKeyFromLines(search_key, lines)
		want := ""

		if got != want {
			t.Errorf("search_key %q, lines %q, got %q, want %q", search_key, lines, got, want)
		}
	})
}