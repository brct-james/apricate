package timecalc_test

import (
	"apricate/timecalc"
	"testing"
	"time"
)

func TestAddSecondsToTimestamp(t *testing.T) {
	t.Run("return empty string when not found", func(t *testing.T) {
		start_time := time.Now()
		seconds := int(10)
		got := timecalc.AddSecondsToTimestamp(start_time, seconds)
		want := start_time.Add(time.Second * time.Duration(seconds))

		if !got.Equal(want) {
			t.Errorf("start_time %q, seconds %q, got %q, want %q", start_time, seconds, got, want)
		}
	})
}