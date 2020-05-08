package middleware

import (
	"net/http"
	"net/url"
)

// Transform Changes request based on query params
// "r": changes Header Referer
// "u": changes Request URL
func Transform(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		if params.Get("r") != "" {
			r.Header["Referer"] = []string{params.Get("r")}
		}

		if params.Get("u") != "" {
			u, err := url.Parse(params.Get("u"))
			if err == nil {
				r.URL = u
				r.RequestURI = u.String()
			}
		}

		h.ServeHTTP(w, r)
	})
}
