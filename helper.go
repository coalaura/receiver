package main

import (
	"net/http"
	"os"
	"regexp"

	"github.com/go-chi/chi/v5"
)

func CleanName(r *http.Request) (string, bool) {
	name := chi.URLParam(r, "name")

	rgx := regexp.MustCompile(`(?m)^[\w.-]+$`)
	if !rgx.MatchString(name) {
		return name, false
	}

	return name, true
}

func EnsureDirectory(path string) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return
	}

	os.MkdirAll(path, 0755)
}

func Respond(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)

	w.Write([]byte(msg))
}
