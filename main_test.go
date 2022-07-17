package main

import (
	"testing"

	"apricate/comparison"
)

func TestHandleServerConfigFlushDBs(t *testing.T) {
	t.Run("return flush_dbs=false when initially true", func(t *testing.T) {
		lines := []string{"flush_dbs=true"}
		got := HandleServerConfigFlushDBs(lines)
		want := []string{"flush_dbs=false"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("return flush_dbs=dev when initially dev", func(t *testing.T) {
		lines := []string{"flush_dbs=dev"}
		got := HandleServerConfigFlushDBs(lines)
		want := []string{"flush_dbs=dev"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("return flush_dbs=false when initially false", func(t *testing.T) {
		lines := []string{"flush_dbs=dev"}
		got := HandleServerConfigFlushDBs(lines)
		want := []string{"flush_dbs=dev"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("return flush_dbs=false appended when flush_dbs missing", func(t *testing.T) {
		lines := []string{"other_key=boring"}
		got := HandleServerConfigFlushDBs(lines)
		want := []string{"other_key=boring", "flush_dbs=false"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("return flush_dbs=false inserted when flush_dbs missing and first line empty", func(t *testing.T) {
		lines := []string{""}
		got := HandleServerConfigFlushDBs(lines)
		want := []string{"flush_dbs=false"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("return flush_dbs=false inserted when flush_dbs missing and first line empty", func(t *testing.T) {
		lines := []string{"", "other_key=boring"}
		got := HandleServerConfigFlushDBs(lines)
		want := []string{"flush_dbs=false", "other_key=boring"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
}

func TestHandleServerConfigRegenerateAuthSecret(t *testing.T) {
	t.Run("return regenerate_auth_secret=false when initially true", func(t *testing.T) {
		lines := []string{"regenerate_auth_secret=true"}
		got := HandleServerConfigRegenerateAuthSecret(lines)
		want := []string{"regenerate_auth_secret=false"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("return regenerate_auth_secret=dev when initially dev", func(t *testing.T) {
		lines := []string{"regenerate_auth_secret=dev"}
		got := HandleServerConfigRegenerateAuthSecret(lines)
		want := []string{"regenerate_auth_secret=dev"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("return regenerate_auth_secret=false when initially false", func(t *testing.T) {
		lines := []string{"regenerate_auth_secret=dev"}
		got := HandleServerConfigRegenerateAuthSecret(lines)
		want := []string{"regenerate_auth_secret=dev"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("return regenerate_auth_secret=false appended when regenerate_auth_secret missing", func(t *testing.T) {
		lines := []string{"other_key=boring"}
		got := HandleServerConfigRegenerateAuthSecret(lines)
		want := []string{"other_key=boring", "regenerate_auth_secret=false"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("return regenerate_auth_secret=false inserted when regenerate_auth_secret missing and first line empty", func(t *testing.T) {
		lines := []string{""}
		got := HandleServerConfigRegenerateAuthSecret(lines)
		want := []string{"regenerate_auth_secret=false"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("return regenerate_auth_secret=false inserted when regenerate_auth_secret missing and first line empty", func(t *testing.T) {
		lines := []string{"", "other_key=boring"}
		got := HandleServerConfigRegenerateAuthSecret(lines)
		want := []string{"regenerate_auth_secret=false", "other_key=boring"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
}

func TestProcessServerConfig(t *testing.T) {
	t.Run("empty config file", func(t *testing.T) {
		lines := []string{""}
		got := ProcessServerConfig(lines)
		want := []string{"flush_dbs=false", "regenerate_auth_secret=false"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("flush=dev regenerate=false", func(t *testing.T) {
		lines := []string{"flush_dbs=dev", "regenerate_auth_secret=false"}
		got := ProcessServerConfig(lines)
		want := []string{"flush_dbs=dev", "regenerate_auth_secret=false"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("flush=dev regenerate=true", func(t *testing.T) {
		lines := []string{"flush_dbs=dev", "regenerate_auth_secret=true"}
		got := ProcessServerConfig(lines)
		want := []string{"flush_dbs=dev", "regenerate_auth_secret=false"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("flush=true regenerate=false", func(t *testing.T) {
		lines := []string{"flush_dbs=true", "regenerate_auth_secret=false"}
		got := ProcessServerConfig(lines)
		want := []string{"flush_dbs=false", "regenerate_auth_secret=false"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("flush=false regenerate=true", func(t *testing.T) {
		lines := []string{"flush_dbs=false", "regenerate_auth_secret=true"}
		got := ProcessServerConfig(lines)
		want := []string{"flush_dbs=false", "regenerate_auth_secret=false"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
	t.Run("flush=true regenerate=true", func(t *testing.T) {
		lines := []string{"flush_dbs=true", "regenerate_auth_secret=true"}
		got := ProcessServerConfig(lines)
		want := []string{"flush_dbs=false", "regenerate_auth_secret=false"}

		if !comparison.StringSlicesEqual(got, want) {
			t.Errorf("lines %q, got %q, want %q", lines, got, want)
		}
	})
}