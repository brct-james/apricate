// Package handlers provides functions for handling web routes
package handlers

import (
	"fmt"
	"net/http"
)

func ServerStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}