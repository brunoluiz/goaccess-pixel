package handler

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"

	"github.com/brunoluiz/goaccess-pixel/middleware"
)

// Middleware function to return middlewares
type Middleware func(h http.Handler) http.Handler

// PixelLogger returns pixel and logs requests to it oin specified output
func PixelLogger(output io.Writer) http.Handler {
	return Pixel(WithTransform(), WithLog(output))
}

// Pixel return pixel to user
func Pixel(middlewares ...Middleware) http.Handler {
	var h http.Handler
	h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
		img.Set(1, 1, color.NRGBA{R: 255, G: 255, B: 255, A: 0})
		png.Encode(w, img)
	})

	for _, m := range middlewares {
		h = m(h)
	}

	return h
}

// WithLog log requests in apache combined format
func WithLog(output io.Writer) Middleware {
	return middleware.Log(output)
}

// WithTransform transform request params to specific headers used by apache logger
func WithTransform() Middleware {
	return middleware.Transform
}
