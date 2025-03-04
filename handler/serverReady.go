package handler

import (
	"net/http"
	"sync/atomic"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
}

func ServerReady(responseWriter http.ResponseWriter, req *http.Request) {

	responseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8") // sets the header of responseWriteer
	responseWriter.WriteHeader(200)
	s := []byte("OK")
	responseWriter.Write(s)
}

func (cfg *ApiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
