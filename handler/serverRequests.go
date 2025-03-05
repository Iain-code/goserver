package handler

import (
	"encoding/json"
	"goserver/internal/database"
	"net/http"
	"strings"
	"sync/atomic"
)

type ApiConfig struct {
	Db             *database.Queries
	fileserverHits atomic.Int32
	Platform       string
}

func ServerReady(responseWriter http.ResponseWriter, req *http.Request) {

	responseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8") // sets the header of responseWriteer
	responseWriter.WriteHeader(200)
	s := []byte("OK")
	responseWriter.Write(s)
}

func Validate(w http.ResponseWriter, r *http.Request) {

	type incPost struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	post := incPost{}
	err := decoder.Decode(&post)

	if err != nil {
		respondWithError(w, 400, "Error decoding parameters")
		return
	}

	if len(post.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleanedText := badWordReplacement(post.Body)

	respondWithJSON(w, 200, cleanedText)

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

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	type responseJSON struct {
		Cleaned_body interface{} `json:"cleaned_body"`
	}

	response := responseJSON{}
	response.Cleaned_body = payload

	data, err := json.Marshal(response)
	if err != nil {
		respondWithError(w, 500, "Error marshalling JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)

}

func badWordReplacement(text string) string {

	words := strings.Split(text, " ")

	for i, word := range words {
		lowered := strings.ToLower(word)
		if lowered == "kerfuffle" || lowered == "sharbert" || lowered == "fornax" {
			words[i] = "****"
		}
	}
	replacement := strings.Join(words, " ")
	return replacement
}
