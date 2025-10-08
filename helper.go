package main

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

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

func ParseAspectRatio(r *http.Request) (float64, error) {
	raw := chi.URLParam(r, "ratio")
	if raw == "" {
		return 0, nil
	}

	index := strings.Index(raw, ":")
	if index == -1 {
		return 0, errors.New("invalid aspect ratio")
	}

	width, err := strconv.ParseInt(raw[:index], 10, 64)
	if err != nil {
		return 0, err
	}

	height, err := strconv.ParseInt(raw[index+1:], 10, 64)
	if err != nil {
		return 0, err
	}

	return float64(width) / float64(height), nil
}

func Respond(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)

	w.Write([]byte(msg))
}
