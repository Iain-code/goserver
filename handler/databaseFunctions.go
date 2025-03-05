package handler

import (
	"database/sql"
	"encoding/json"
	"goserver/internal/database"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) NewUser(w http.ResponseWriter, r *http.Request) {

	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	type Email struct {
		Email string `json:"email"`
	}
	email := Email{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&email)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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

	user.Email = email.Email

	newUser, err := apiCfg.Db.CreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userStruct := User{}
	userStruct.ID = newUser.ID
	userStruct.Email = newUser.Email
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
