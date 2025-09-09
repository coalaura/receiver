package main

import (
	"net/http"

	"github.com/coalaura/plain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var log = plain.New(plain.WithDate(plain.RFC3339Local))

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(log.Middleware())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		Respond(w, 200, "## `POST /file/{name}`\nStores file from `file` post field in `files/{name}`\n\n## `POST /image/{name}`\nStores image from `file` post field as webp in `files/{name}.webp`\n\n## Response\nEndpoints respond with `text/plain`.\n- Success: status=200, content=`OK`\n- Fail: status!=200, content=`ERROR`")
	})

	r.Post("/file/{name}", HandleFileUpload)
	r.Post("/image/{name}", HandleImageUpload)

	log.Println("Listening at http://localhost:6942/")
	http.ListenAndServe(":6942", r)
}
