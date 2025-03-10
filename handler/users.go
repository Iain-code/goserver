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

	email := received{}

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
	userStruct := &User{}
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

func (apiCfg *ApiConfig) Login(w http.ResponseWriter, r *http.Request) {

	type TokenUser struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
	}

	receivedData := received{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&receivedData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// FIND THE USERS DATA USING THE EMAIL IN JSON REQUEST

	data, err := apiCfg.Db.FindUserEmail(r.Context(), receivedData.Email)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	err = auth.CheckPasswordHash(receivedData.Password, data.HashedPassword.String)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "refresh token failed")
		return
	}
	RefreshTokenParams := database.MakeRefreshTokenParams{}
	RefreshTokenParams.Token = refreshToken
	RefreshTokenParams.ExpiresAt.Time = time.Now().Add(60 * 24 * time.Hour)
	RefreshTokenParams.ExpiresAt.Valid = true
	RefreshTokenParams.CreatedAt = time.Now()
	RefreshTokenParams.UpdatedAt = time.Now()
	RefreshTokenParams.UserID = data.ID
	err = apiCfg.Db.MakeRefreshToken(r.Context(), RefreshTokenParams)
	if err != nil {
		respondWithError(w, 500, "error making refresh token")
		return
	}

	tokenstr, err := auth.MakeJWT(data.ID, apiCfg.TokenSecret, time.Hour*1)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	ID, err := auth.ValidateJWT(tokenstr, apiCfg.TokenSecret)
	if ID != data.ID {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	refreshTokenStruct := TokenUser{}
	refreshTokenStruct.ID = data.ID
	refreshTokenStruct.CreatedAt = RefreshTokenParams.CreatedAt
	refreshTokenStruct.UpdatedAt = RefreshTokenParams.UpdatedAt
	refreshTokenStruct.Email = data.Email
	refreshTokenStruct.Token = tokenstr
	refreshTokenStruct.RefreshToken = refreshToken
	refreshTokenStruct.IsChirpyRed = data.IsChirpyRed

	respondWithJSON(w, 200, refreshTokenStruct)
}

func (ApiCfg *ApiConfig) UpdateDetails(w http.ResponseWriter, r *http.Request) {

	type UpdatePassEmail struct {
		ID             uuid.UUID      `json:"id"`
		Email          string         `json:"email"`
		HashedPassword sql.NullString `json:"-"`
	}

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

	request := received{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hash, err := auth.HashPassword(request.Password)
	if err != nil {
		log.Printf("error while hashing password for user %v", request.Email)
		return
	}

	updated := database.UpdatePassEmailParams{}
	updated.ID = number
	updated.Email = request.Email
	updated.HashedPassword = sql.NullString{
		String: hash,
		Valid:  true,
	}
	err = ApiCfg.Db.UpdatePassEmail(r.Context(), updated)

	rtn := UpdatePassEmail{}
	rtn.Email = updated.Email
	rtn.HashedPassword = updated.HashedPassword
	rtn.ID = updated.ID

	if err != nil {
		respondWithError(w, 400, "failed to update user email or password")
		return
	}

	respondWithJSON(w, 200, rtn)
}

func (apiCfg *ApiConfig) AddChirpyRed(w http.ResponseWriter, r *http.Request) {

	type Event struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	event := Event{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&event)
	if err != nil {
		respondWithError(w, 401, "Invalid request body")
	}
	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "Invalid request body")
		return
	}

	if key != apiCfg.PolkaKey {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	if event.Event != "user.upgraded" {
		respondWithError(w, 204, "")
		return
	}
	num, err := uuid.Parse(event.Data.UserID)

	if err != nil {
		respondWithError(w, 400, "Bad Request")
		return
	}

	err = apiCfg.Db.AddChirpyRed(r.Context(), num)
	if err != nil {
		respondWithError(w, 404, "Not Found")
		return
	}
	respondWithJSON(w, 204, "")
}
