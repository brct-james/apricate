package log_test

import (
	"apricate/log"
	"testing"
)

// Formatting functions

func TestBold(t *testing.T) {
	start := "String"
	got := log.Bold(start)
	want := "\u001b[1mString\u001b[0m"

	if got != want {
		t.Errorf("start %s got %s want %s", start, got, want)
	}
}

// Coloring functions

func TestBlue(t *testing.T) {
	start := "String"
	got := log.Blue(start)
	want := "\u001b[34mString\u001b[0m"

	if got != want {
		t.Errorf("start %s got %s want %s", start, got, want)
	}
}

func TestYellow(t *testing.T) {
	start := "String"
	got := log.Yellow(start)
	want := "\u001b[33mString\u001b[0m"

	if got != want {
		t.Errorf("start %s got %s want %s", start, got, want)
	}
}

func TestRed(t *testing.T) {
	start := "String"
	got := log.Red(start)
	want := "\u001b[31mString\u001b[0m"

	if got != want {
		t.Errorf("start %s got %s want %s", start, got, want)
	}
}

func TestCyan(t *testing.T) {
	start := "String"
	got := log.Cyan(start)
	want := "\u001b[36mString\u001b[0m"

	if got != want {
		t.Errorf("start %s got %s want %s", start, got, want)
	}
}

func TestWhite(t *testing.T) {
	start := "String"
	got := log.White(start)
	want := "\u001b[37mString\u001b[0m"

	if got != want {
		t.Errorf("start %s got %s want %s", start, got, want)
	}
}

func TestGreen(t *testing.T) {
	start := "String"
	got := log.Green(start)
	want := "\u001b[32mString\u001b[0m"

	if got != want {
		t.Errorf("start %s got %s want %s", start, got, want)
	}
}

// Background functions

func TestMagentaBackground(t *testing.T) {
	start := "String"
	got := log.MagentaBackground(start)
	want := "\u001b[45mString\u001b[0m"

	if got != want {
		t.Errorf("start %s got %s want %s", start, got, want)
	}
}

func TestCyanBackground(t *testing.T) {
	start := "String"
	got := log.CyanBackground(start)
	want := "\u001b[46mString\u001b[0m"

	if got != want {
		t.Errorf("start %s got %s want %s", start, got, want)
	}
}

// Test Pass Fail functions

func TestPass(t *testing.T) {
	start := "String"
	got := log.Pass(start)
	want := "\u001b[42m\u001b[39mString\u001b[0m"

	if got != want {
		t.Errorf("start %s got %s want %s", start, got, want)
	}
}

func TestFail(t *testing.T) {
	start := "String"
	got := log.Fail(start)
	want := "\u001b[41m\u001b[39mString\u001b[0m"

	if got != want {
		t.Errorf("start %s got %s want %s", start, got, want)
	}
}

func TestFormatPassFail(t *testing.T) {
	t.Run("format green when pass", func(t *testing.T) {
		start := "String"
		pass_string := "String"
		got := log.FormatPassFail(start, pass_string)
		want := "\u001b[42m\u001b[39mString\u001b[0m"

		if got != want {
			t.Errorf("start %s pass_string %s got %s want %s", start, pass_string, got, want)
		}
	})
	t.Run("format red when fail", func(t *testing.T) {
		start := "String"
		pass_string := "NotString"
		got := log.FormatPassFail(start, pass_string)
		want := "\u001b[41m\u001b[39mString\u001b[0m"

		if got != want {
			t.Errorf("start %s pass_string %s got %s want %s", start, pass_string, got, want)
		}
	})
}