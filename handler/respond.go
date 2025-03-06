package handler

import (
	"encoding/json"
	"net/http"
)

func ServerReady(responseWriter http.ResponseWriter, req *http.Request) {

	responseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8") // sets the header of responseWriteer
	responseWriter.WriteHeader(200)
	s := []byte("OK")
	responseWriter.Write(s)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	data, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 500, "Error marshalling JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)

}

func respondWithError(w http.ResponseWriter, code int, msg string) {

	type ErrorResponse struct {
		Error string `json:"error"`
	}

	errResp := ErrorResponse{}
	errResp.Error = msg
	data, _ := json.Marshal(errResp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
