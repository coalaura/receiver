package main

import (
	_ "embed"
	"net/http"

	"github.com/coalaura/plain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	log    = plain.New(plain.WithDate(plain.RFC3339Local))
	worker = NewWorker(8)

	//go:embed README.md
	readme string
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(log.Middleware())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		Respond(w, 200, readme)
	})

	r.Post("/file/{name}", HandleFileUpload)

	r.Post("/image/{name}", HandleImageUpload)
	r.Post("/image/{name}/{size}", HandleImageUpload)
	r.Post("/image/{name}/{size}/{ratio}", HandleImageUpload)

	log.Println("Listening at http://localhost:6942/")
	http.ListenAndServe(":6942", r)
}
