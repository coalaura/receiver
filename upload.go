package main

import (
	"bytes"
	"image"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"image/draw"
	_ "image/jpeg"
	_ "image/png"

	"github.com/gen2brain/webp"
	"github.com/go-chi/chi/v5"
	"github.com/nfnt/resize"
)

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

	go func() {
		log.Printf("Uploading %q\n", name)

		file, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			log.Warnf("Failed to open target file: %v\n", err)

			return
		}

		defer file.Close()

		_, err = io.Copy(file, buf)
		if err != nil {
			log.Warnf("Failed to copy to target file: %v\n", err)

			return
		}

		log.Printf("Finished uploading %q\n", name)
	}()

	Respond(w, http.StatusOK, "OK")
}

func HandleImageUpload(w http.ResponseWriter, r *http.Request) {
	name, ok := CleanName(r)
	if !ok {
		log.Warnf("Invalid file name: %q\n", name)

		Respond(w, http.StatusBadRequest, "ERROR")

		return
	}

	var size uint

	if raw := chi.URLParam(r, "size"); raw != "" {
		num, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			log.Warnf("Invalid size parameter: %v\n", err)

			Respond(w, http.StatusBadRequest, "ERROR")

			return
		}

		if num != 0 {
			size = min(max(uint(num), 128), 4096)
		}
	}

	ratio, err := ParseAspectRatio(r)
	if err != nil {
		log.Warnf("Invalid ratio parameter: %v\n", err)

		Respond(w, http.StatusBadRequest, "ERROR")

		return
	}

	buf, err := ReceiveFile(r)
	if err != nil {
		log.Warnf("Failed to receive file: %v\n", err)

		Respond(w, http.StatusBadRequest, "ERROR")

		return
	}

	go func() {
		log.Printf("Uploading %q\n", name)

		img, _, err := image.Decode(buf)
		if err != nil {
			log.Warnf("Failed to receive image: %v\n", err)

			return
		}

		if index := strings.LastIndex(name, "."); index != -1 {
			name = name[:index]
		}

		name += ".webp"

		if ratio > 0 {
			log.Printf("Cropping %q to ratio %.2f\n", name, ratio)

			img = CropToRatio(img, ratio)
		}

		if size > 0 {
			log.Printf("Resizing %q to max %d\n", name, size)

			img = resize.Thumbnail(size, size, img, resize.Lanczos3)
		}

		file, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			log.Warnf("Failed to open target file: %v\n", err)

			Respond(w, http.StatusInternalServerError, "ERROR")

			return
		}

		defer file.Close()

		err = webp.Encode(file, img, webp.Options{
			Quality: 90,
			Method:  6,
		})
		if err != nil {
			log.Warnf("Failed to encode webp: %v\n", err)

			Respond(w, http.StatusInternalServerError, "ERROR")

			return
		}

		log.Printf("Finished uploading %q\n", name)
	}()

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

func CropToRatio(img image.Image, ratio float64) image.Image {
	bounds := img.Bounds()

	width := bounds.Dx()
	height := bounds.Dy()

	if width <= 0 || height <= 0 {
		return img
	}

	current := float64(width) / float64(height)

	if current == ratio {
		return img
	}

	nWidth := width
	nHeight := height

	if current > ratio {
		nWidth = int(math.Round(float64(height) * ratio))

		if nWidth < 1 {
			nWidth = 1
		}
	} else if current < ratio {
		nHeight = int(math.Round(float64(width) / ratio))

		if nHeight < 1 {
			nHeight = 1
		}
	}

	x0 := bounds.Min.X + (width-nWidth)/2
	y0 := bounds.Min.Y + (height-nHeight)/2

	crop := image.Rect(x0, y0, x0+nWidth, y0+nHeight)

	dst := image.NewRGBA(image.Rect(0, 0, nWidth, nHeight))

	draw.Draw(dst, dst.Bounds(), img, crop.Min, draw.Src)

	return dst
}
