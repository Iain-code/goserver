package handler

import (
	"database/sql"
	"encoding/json"
	"goserver/internal/auth"
	"goserver/internal/database"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) NewUser(w http.ResponseWriter, r *http.Request) {

	type User struct {
		ID             uuid.UUID      `json:"id"`
		CreatedAt      time.Time      `json:"created_at"`
		UpdatedAt      time.Time      `json:"updated_at"`
		Email          string         `json:"email"`
		HashedPassword sql.NullString `json:"hashed_password"`
	}

	type Email struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	email := Email{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&email)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	hash, err := auth.HashPassword(email.Password)
	if err != nil {
		log.Printf("error while hashing password for user %v", email.Email)
		return
	}
	user := database.CreateUserParams{}
	user.ID = uuid.New()
	user.CreatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	user.UpdatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	user.HashedPassword = sql.NullString{
		String: hash,
		Valid:  true,
	}
	user.Email = email.Email

	newUser, err := apiCfg.Db.CreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userStruct := User{}
	userStruct.ID = newUser.ID
	userStruct.Email = newUser.Email
	userStruct.HashedPassword = newUser.HashedPassword
	// If using sql.NullTime
	if newUser.CreatedAt.Valid {
		userStruct.CreatedAt = newUser.CreatedAt.Time
	}
	if newUser.UpdatedAt.Valid {
		userStruct.UpdatedAt = newUser.UpdatedAt.Time
	}

	data, err := json.Marshal(userStruct)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(data)
}

func (apiCfg *ApiConfig) Login(w http.ResponseWriter, r *http.Request) {

	type User struct {
		ID             uuid.UUID      `json:"id"`
		CreatedAt      time.Time      `json:"created_at"`
		UpdatedAt      time.Time      `json:"updated_at"`
		Email          string         `json:"email"`
		HashedPassword sql.NullString `json:"-"`
		Token          string         `json:"token"`
	}

	type received struct {
		Password           string        `json:"password"`
		Email              string        `json:"email"`
		Expires_in_seconds time.Duration `json:"optional_field,expires_in_seconds"`
	}

	user := User{}
	receivedData := received{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&receivedData)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	data, err := apiCfg.Db.FindUserEmail(r.Context(), receivedData.Email)
	if err != nil {
		respondWithError(w, 401, "Unauthorized 1")
		return
	}
	user.ID = data.ID
	if data.CreatedAt.Valid {
		user.CreatedAt = data.CreatedAt.Time
	}
	if data.UpdatedAt.Valid {
		user.UpdatedAt = data.UpdatedAt.Time
	}
	user.Email = data.Email

	err = auth.CheckPasswordHash(receivedData.Password, data.HashedPassword.String)
	if err != nil {
		respondWithError(w, 401, "Unauthorized 2")
		return
	}

	if receivedData.Expires_in_seconds > time.Duration(1)*time.Hour || receivedData.Expires_in_seconds == 0 {
		receivedData.Expires_in_seconds = time.Duration(1) * time.Hour
	}

	tokenstr, err := auth.MakeJWT(user.ID, apiCfg.TokenSecret, receivedData.Expires_in_seconds)
	if err != nil {
		respondWithError(w, 401, "Unauthorized 3")
		return
	}
	ID, err := auth.ValidateJWT(tokenstr, apiCfg.TokenSecret)
	if ID != user.ID {
		respondWithError(w, 401, "Unauthorized 4")
		return
	}
	user.Token = tokenstr
	respondWithJSON(w, 200, user)
}
