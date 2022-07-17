package filemngr_test

import (
	"apricate/filemngr"
	"testing"
)

func TestGetKeyFromLines(t *testing.T) {
	t.Run("return empty string when not found", func(t *testing.T) {
		search_key := "key"
		lines := []string{"not key value", "incorrect_key=value"}
		got_i, got := filemngr.GetKeyFromLines(search_key, lines)
		want := ""
		want_i := -1

		if got != want || got_i != want_i {
			t.Errorf("search_key %q, lines %q, got %q, want %q, got_i %q, want_i %q", search_key, lines, got, want, got_i, want_i)
		}
	})
	t.Run("return empty string when lines empty", func(t *testing.T) {
		search_key := "key"
		lines := []string{}
		got_i, got := filemngr.GetKeyFromLines(search_key, lines)
		want := ""
		want_i := -1

		if got != want || got_i != want_i {
			t.Errorf("search_key %q, lines %q, got %q, want %q, got_i %q, want_i %q", search_key, lines, got, want, got_i, want_i)
		}
	})
	t.Run("return value string when found", func(t *testing.T) {
		search_key := "key"
		lines := []string{"not key value", "key=value"}
		got_i, got := filemngr.GetKeyFromLines(search_key, lines)
		want := "value"
		want_i := 1

		if got != want || got_i != want_i {
			t.Errorf("search_key %q, lines %q, got %q, want %q, got_i %q, want_i %q", search_key, lines, got, want, got_i, want_i)
		}
	})
	t.Run("return value string when found and value contains =", func(t *testing.T) {
		search_key := "key"
		lines := []string{"key=value=test"}
		got_i, got := filemngr.GetKeyFromLines(search_key, lines)
		want := "value=test"
		want_i := 0

		if got != want || got_i != want_i {
			t.Errorf("search_key %q, lines %q, got %q, want %q, got_i %q, want_i %q", search_key, lines, got, want, got_i, want_i)
		}
	})
	t.Run("return empty string for notkey=key=value", func(t *testing.T) {
		search_key := "key"
		lines := []string{"notkey=key=value"}
		got_i, got := filemngr.GetKeyFromLines(search_key, lines)
		want := ""
		want_i := -1

		if got != want || got_i != want_i {
			t.Errorf("search_key %q, lines %q, got %q, want %q, got_i %q, want_i %q", search_key, lines, got, want, got_i, want_i)
		}
	})
}