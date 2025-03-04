package handler

import (
	"net/http"
)

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do something before the next handler
		cfg.fileserverHits.Add(1)

		// Call the next handler to continue processing
		next.ServeHTTP(w, r)
	})
}
