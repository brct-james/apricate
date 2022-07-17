package main

import (
	"testing"

	"apricate/comparison"
)

func TestHandleServerConfigFlushDBs(t *testing.T) {
	t.Run("return flush_dbs=false when initially true", func(t *testing.T) {
		lines := []string{"flush_dbs=true"}
		got := ProcessServerConfig(lines)
		want := []string{"flush_dbs=false"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("return flush_dbs=dev when initially dev", func(t *testing.T) {
		lines := []string{"flush_dbs=dev"}
		got := ProcessServerConfig(lines)
		want := []string{"flush_dbs=dev"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("return flush_dbs=false when initially false", func(t *testing.T) {
		lines := []string{"flush_dbs=dev"}
		got := ProcessServerConfig(lines)
		want := []string{"flush_dbs=dev"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("return flush_dbs=false appended when flush_dbs missing", func(t *testing.T) {
		lines := []string{"other_key=boring"}
		got := ProcessServerConfig(lines)
		want := []string{"other_key=boring", "flush_dbs=false"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("return flush_dbs=false inserted when flush_dbs missing and first line empty", func(t *testing.T) {
		lines := []string{""}
		got := ProcessServerConfig(lines)
		want := []string{"flush_dbs=false"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("return flush_dbs=false inserted when flush_dbs missing and first line empty", func(t *testing.T) {
		lines := []string{"", "other_key=boring"}
		got := ProcessServerConfig(lines)
		want := []string{"flush_dbs=false", "other_key=boring"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
}