package middleware

import (
	"io"
	"net/http"

	"github.com/gorilla/handlers"
)

// Log logs output to apache combined format
func Log(output io.Writer) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return (handlers.CombinedLoggingHandler(output, h))
	}
}
