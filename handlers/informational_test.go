package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETServerStatus(t *testing.T) {
	t.Run("returns server status", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/status", nil)
		response := httptest.NewRecorder()

		ServerStatus(response, request)

		got := response.Body.String()
		want := "OK"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}