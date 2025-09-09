package main

import (
	"bytes"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gen2brain/webp"
)

const ReceiveDirectory = "files"

func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	name, ok := CleanName(r)
	if !ok {
		log.Warnf("Invalid file name: %q\n", name)

		Respond(w, http.StatusBadRequest, "ERROR")

		return
	}

	buf, err := ReceiveFile(r)
	if err != nil {
		log.Warnf("Failed to receive file: %v\n", err)

		Respond(w, http.StatusBadRequest, "ERROR")

		return
	}

	EnsureDirectory(ReceiveDirectory)

	path := filepath.Join(ReceiveDirectory, name)

	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Warnf("Failed to open target file: %v\n", err)

		Respond(w, http.StatusInternalServerError, "ERROR")

		return
	}

	defer file.Close()

	_, err = io.Copy(file, buf)
	if err != nil {
		log.Warnf("Failed to copy to target file: %v\n", err)

		Respond(w, http.StatusInternalServerError, "ERROR")

		return
	}

	Respond(w, http.StatusOK, "OK")
}

func HandleImageUpload(w http.ResponseWriter, r *http.Request) {
	name, ok := CleanName(r)
	if !ok {
		log.Warnf("Invalid file name: %q\n", name)

		Respond(w, http.StatusBadRequest, "ERROR")

		return
	}

	img, err := ReceiveImage(r)
	if err != nil {
		log.Warnf("Failed to receive image: %v\n", err)

		Respond(w, http.StatusBadRequest, "ERROR")

		return
	}

	EnsureDirectory(ReceiveDirectory)

	path := filepath.Join(ReceiveDirectory, name)

	if index := strings.LastIndex(path, "."); index != -1 {
		path = path[:index]
	}

	path += ".webp"

	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Warnf("Failed to open target file: %v\n", err)

		Respond(w, http.StatusInternalServerError, "ERROR")

		return
	}

	defer file.Close()

	err = webp.Encode(file, img, webp.Options{
		Quality: 90,
		Method:  5,
	})
	if err != nil {
		log.Warnf("Failed to encode webp: %v\n", err)

		Respond(w, http.StatusInternalServerError, "ERROR")

		return
	}

	Respond(w, http.StatusOK, "OK")
}

func ReceiveFile(r *http.Request) (io.Reader, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, err
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var buf bytes.Buffer

	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

func ReceiveImage(r *http.Request) (image.Image, error) {
	buf, err := ReceiveFile(r)
	if err != nil {
		return nil, err
	}

	input, _, err := image.Decode(buf)
	if err != nil {
		return nil, err
	}

	return input, nil
}
