package route

import "net/http"

// unsecured is a SecurityHandler for explicitly acknowledging that a route is wide open for use.
func unsecured() SecurityHandler {
	return func(h http.Handler) http.Handler {
		return h
	}
}

// read https://developer.mozilla.org/en-US/docs/Web/HTTP/Content_negotiation#Server-driven_negotiation
func contentTypeSecurity(rule string) SecurityHandler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if contentType := r.Header.Get("Content-Type"); contentType != rule {
				w.WriteHeader(http.StatusNotAcceptable)
				w.Write([]byte("Content type is wrong."))
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
