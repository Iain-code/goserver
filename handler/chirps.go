package handler

import (
	"encoding/json"
	"goserver/internal/auth"
	"goserver/internal/database"
	"net/http"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

type ApiConfig struct {
	Db             *database.Queries
	fileserverHits atomic.Int32
	Platform       string
	TokenSecret    string
	PolkaKey       string
}

func (apiCfg *ApiConfig) Chirps(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	post := incPost{}
	err := decoder.Decode(&post)

	if err != nil {
		respondWithError(w, 400, "error decoding parameters")
		return
	}

	if len(post.Body) > 140 {
		respondWithError(w, 400, "chirp is too long")
		return
	}

	headers := r.Header
	token, err := auth.GetBearerToken(headers)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	ID, err := auth.ValidateJWT(token, apiCfg.TokenSecret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	cleanedText := badWordReplacement(post.Body)

	chirps := database.CreateChirpParams{}
	chirps.ID = uuid.New()
	chirps.CreatedAt = time.Now()
	chirps.UpdatedAt = time.Now()
	chirps.Body = cleanedText
	chirps.UserID = uuid.NullUUID{
		UUID:  ID,
		Valid: true,
	}

	data, err := apiCfg.Db.CreateChirp(r.Context(), chirps)
	if err != nil {
		respondWithError(w, 400, "failed to create chirp")
		return
	}
	chirp := CreateChirp{}
	chirp.ID = data.ID
	chirp.CreatedAt = data.CreatedAt
	chirp.UpdatedAt = data.UpdatedAt
	chirp.Body = data.Body
	chirp.UserID = data.UserID

	respondWithJSON(w, 201, chirp)

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

func (apiCfg *ApiConfig) GetAllChirps(w http.ResponseWriter, r *http.Request) {

	type CreateChirp struct {
		ID        uuid.UUID     `json:"id"`
		CreatedAt time.Time     `json:"created_at"`
		UpdatedAt time.Time     `json:"updated_at"`
		Body      string        `json:"body"`
		UserID    uuid.NullUUID `json:"user_id"`
	}

	chirpSlice := []CreateChirp{}
	s := r.URL.Query().Get("author_id")
	chirps := []database.Chirp{}

	if s == "" {
		var err error
		chirps, err = apiCfg.Db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, 500, "get all chirps failed")
			return
		}
	} else {

		var authorID uuid.NullUUID
		num, err := uuid.Parse(s)
		authorID.UUID = num
		authorID.Valid = true

		if err != nil {
			respondWithError(w, 400, "Bad Request")
			return
		}

		chirps, err = apiCfg.Db.GetAllChirpsFromUser(r.Context(), authorID)

		if err != nil {
			respondWithError(w, 500, "get all chirps failed")
			return
		}
	}
	for _, chirp := range chirps {
		newChirp := CreateChirp{}
		newChirp.ID = chirp.ID
		newChirp.CreatedAt = chirp.CreatedAt
		newChirp.UpdatedAt = chirp.UpdatedAt
		newChirp.Body = chirp.Body
		newChirp.UserID = chirp.UserID
		chirpSlice = append(chirpSlice, newChirp)
	}

	sortQuery := r.URL.Query().Get("sort")

	if sortQuery == "asc" || sortQuery == "" {
		sort.Slice(chirpSlice, func(i, j int) bool {
			return chirpSlice[i].CreatedAt.Before(chirpSlice[j].CreatedAt)
		})
	} else {

		if sortQuery == "desc" {
			sort.Slice(chirpSlice, func(i, j int) bool {
				return chirpSlice[i].CreatedAt.After(chirpSlice[j].CreatedAt)
			})
		}
	}

	data, err := json.Marshal(chirpSlice)
	if err != nil {
		respondWithError(w, 500, "Error marshalling JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)

}

func (apiCfg *ApiConfig) GetOneChirp(w http.ResponseWriter, r *http.Request) {

	type CreateChirp struct {
		ID        uuid.UUID     `json:"id"`
		CreatedAt time.Time     `json:"created_at"`
		UpdatedAt time.Time     `json:"updated_at"`
		Body      string        `json:"body"`
		UserID    uuid.NullUUID `json:"user_id"`
	}

	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr) // changes a string into a UUID
	if err != nil {
		// Handle invalid UUID format
		http.Error(w, "Invalid chirp ID format", http.StatusBadRequest)
		return
	}

	chirp, err := apiCfg.Db.GetOneChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "error getting one chirp")
		return
	}
	newChirp := CreateChirp{}
	newChirp.ID = chirp.ID
	newChirp.CreatedAt = chirp.CreatedAt
	newChirp.UpdatedAt = chirp.UpdatedAt
	newChirp.Body = chirp.Body
	newChirp.UserID = chirp.UserID
	respondWithJSON(w, 200, newChirp)

}

func (ApiCfg *ApiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {

	tkn, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "failed to get bearer token")
		return
	}

	number, err := auth.ValidateJWT(tkn, ApiCfg.TokenSecret)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr) // changes a string into a UUID
	if err != nil {
		// Handle invalid UUID format
		http.Error(w, "Invalid chirp ID format", http.StatusBadRequest)
		return
	}
	chirp, err := ApiCfg.Db.GetOneChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "error getting one chirp")
		return
	}

	// If chirp.UserID is a uuid.NullUUID and authenticatedUserID is a uuid.UUID
	if !chirp.UserID.Valid || chirp.UserID.UUID != number {
		// Return 403 Forbidden - either the user ID is not valid or not matching
		respondWithError(w, 403, "userID not valid")
		return
	}
	err = ApiCfg.Db.DeleteChirp(r.Context(), chirp.ID)
	if err != nil {
		respondWithError(w, 404, "chirp not found")
		return
	}

	respondWithJSON(w, 204, "Chirp Deleted")

}
