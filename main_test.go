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

func TestGetValueFromServerConfigByKey(t *testing.T) {
	t.Run("get when key present", func(t *testing.T) {
		lines := []string{"test_key=true"}
		key := "test_key"
		default_value := "false"
		got_index, got_value, got_lines := GetValueFromServerConfigByKey(lines, key, default_value)
		want_index := 0
		want_value := "true"
		want_lines := []string{"test_key=true"}

		if !comparison.StringSlicesEqual(got_lines, want_lines) || got_index != want_index || got_value != want_value {
			t.Errorf("lines %q, key %q, default_value %q, got_index %v, got_value %q, got_lines %q, want_index %v, want_value %q, want_lines %q", lines, key, default_value, got_index, got_value, got_lines, want_index, want_value, want_lines)
		}
	})
	t.Run("get when lines empty", func(t *testing.T) {
		lines := []string{""}
		key := "test_key"
		default_value := "false"
		got_index, got_value, got_lines := GetValueFromServerConfigByKey(lines, key, default_value)
		want_index := 0
		want_value := "false"
		want_lines := []string{"test_key=false"}

		if !comparison.StringSlicesEqual(got_lines, want_lines) || got_index != want_index || got_value != want_value {
			t.Errorf("lines %q, key %q, default_value %q, got_index %v, got_value %q, got_lines %q, want_index %v, want_value %q, want_lines %q", lines, key, default_value, got_index, got_value, got_lines, want_index, want_value, want_lines)
		}
	})
	t.Run("get when key missing", func(t *testing.T) {
		lines := []string{"other_test_key=true"}
		key := "test_key"
		default_value := "true"
		got_index, got_value, got_lines := GetValueFromServerConfigByKey(lines, key, default_value)
		want_index := 1
		want_value := "true"
		want_lines := []string{"other_test_key=true", "test_key=true"}

		if !comparison.StringSlicesEqual(got_lines, want_lines) || got_index != want_index || got_value != want_value {
			t.Errorf("lines %q, key %q, default_value %q, got_index %v, got_value %q, got_lines %q, want_index %v, want_value %q, want_lines %q", lines, key, default_value, got_index, got_value, got_lines, want_index, want_value, want_lines)
		}
	})
	t.Run("get when key missing and blank first line", func(t *testing.T) {
		lines := []string{"", "other_test_key=true"}
		key := "test_key"
		default_value := "true"
		got_index, got_value, got_lines := GetValueFromServerConfigByKey(lines, key, default_value)
		want_index := 0
		want_value := "true"
		want_lines := []string{"test_key=true", "other_test_key=true"}

		if !comparison.StringSlicesEqual(got_lines, want_lines) || got_index != want_index || got_value != want_value {
			t.Errorf("lines %q, key %q, default_value %q, got_index %v, got_value %q, got_lines %q, want_index %v, want_value %q, want_lines %q", lines, key, default_value, got_index, got_value, got_lines, want_index, want_value, want_lines)
		}
	})
}

func TestProcessServerConfig(t *testing.T) {
	t.Run("empty config file", func(t *testing.T) {
		lines := []string{""}
		got_lines, got_misc_config := ProcessServerConfig(lines)
		want_lines := []string{"flush_dbs=false", "regenerate_auth_secret=false", "listen_port=:8080", "redis_addr=rdb:6379", "api_version=0.5.0"}
		want_misc_config := []string{":8080", "rdb:6379", "0.5.0"}

		if !comparison.StringSlicesEqual(got_lines, want_lines) || !comparison.StringSlicesEqual(got_misc_config, want_misc_config) {
			t.Errorf("lines %q, got_lines %q, got_misc_config %q, want_lines %q, want_misc_config %q", lines, got_lines, got_misc_config, want_lines, want_misc_config)
		}
	})
	t.Run("flush=dev regenerate=false", func(t *testing.T) {
		lines := []string{"flush_dbs=dev", "regenerate_auth_secret=false"}
		got_lines, got_misc_config := ProcessServerConfig(lines)
		want_lines := []string{"flush_dbs=dev", "regenerate_auth_secret=false", "listen_port=:8080", "redis_addr=rdb:6379", "api_version=0.5.0"}
		want_misc_config := []string{":8080", "rdb:6379", "0.5.0"}

		if !comparison.StringSlicesEqual(got_lines, want_lines) || !comparison.StringSlicesEqual(got_misc_config, want_misc_config) {
			t.Errorf("lines %q, got_lines %q, got_misc_config %q, want_lines %q, want_misc_config %q", lines, got_lines, got_misc_config, want_lines, want_misc_config)
		}
	})
	t.Run("flush=dev regenerate=true", func(t *testing.T) {
		lines := []string{"flush_dbs=dev", "regenerate_auth_secret=true"}
		got_lines, got_misc_config := ProcessServerConfig(lines)
		want_lines := []string{"flush_dbs=dev", "regenerate_auth_secret=false", "listen_port=:8080", "redis_addr=rdb:6379", "api_version=0.5.0"}
		want_misc_config := []string{":8080", "rdb:6379", "0.5.0"}

		if !comparison.StringSlicesEqual(got_lines, want_lines) || !comparison.StringSlicesEqual(got_misc_config, want_misc_config) {
			t.Errorf("lines %q, got_lines %q, got_misc_config %q, want_lines %q, want_misc_config %q", lines, got_lines, got_misc_config, want_lines, want_misc_config)
		}
	})
	t.Run("flush=true regenerate=false", func(t *testing.T) {
		lines := []string{"flush_dbs=true", "regenerate_auth_secret=false"}
		got_lines, got_misc_config := ProcessServerConfig(lines)
		want_lines := []string{"flush_dbs=false", "regenerate_auth_secret=false", "listen_port=:8080", "redis_addr=rdb:6379", "api_version=0.5.0"}
		want_misc_config := []string{":8080", "rdb:6379", "0.5.0"}

		if !comparison.StringSlicesEqual(got_lines, want_lines) || !comparison.StringSlicesEqual(got_misc_config, want_misc_config) {
			t.Errorf("lines %q, got_lines %q, got_misc_config %q, want_lines %q, want_misc_config %q", lines, got_lines, got_misc_config, want_lines, want_misc_config)
		}
	})
	t.Run("flush=false regenerate=true", func(t *testing.T) {
		lines := []string{"flush_dbs=false", "regenerate_auth_secret=true"}
		got_lines, got_misc_config := ProcessServerConfig(lines)
		want_lines := []string{"flush_dbs=false", "regenerate_auth_secret=false", "listen_port=:8080", "redis_addr=rdb:6379", "api_version=0.5.0"}
		want_misc_config := []string{":8080", "rdb:6379", "0.5.0"}

		if !comparison.StringSlicesEqual(got_lines, want_lines) || !comparison.StringSlicesEqual(got_misc_config, want_misc_config) {
			t.Errorf("lines %q, got_lines %q, got_misc_config %q, want_lines %q, want_misc_config %q", lines, got_lines, got_misc_config, want_lines, want_misc_config)
		}
	})
	t.Run("flush=true regenerate=true", func(t *testing.T) {
		lines := []string{"flush_dbs=true", "regenerate_auth_secret=true"}
		got_lines, got_misc_config := ProcessServerConfig(lines)
		want_lines := []string{"flush_dbs=false", "regenerate_auth_secret=false", "listen_port=:8080", "redis_addr=rdb:6379", "api_version=0.5.0"}
		want_misc_config := []string{":8080", "rdb:6379", "0.5.0"}

		if !comparison.StringSlicesEqual(got_lines, want_lines) || !comparison.StringSlicesEqual(got_misc_config, want_misc_config) {
			t.Errorf("lines %q, got_lines %q, got_misc_config %q, want_lines %q, want_misc_config %q", lines, got_lines, got_misc_config, want_lines, want_misc_config)
		}
	})
}