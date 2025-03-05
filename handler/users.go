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
