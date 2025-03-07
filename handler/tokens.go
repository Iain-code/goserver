package handler

import (
	"database/sql"
	"goserver/internal/auth"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenJSON struct {
	Token     string       `json:"token"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	UserID    uuid.UUID    `json:"user_id"`
	ExpiresAt sql.NullTime `json:"expires_at"`
	RevokedAt sql.NullTime `json:"revoked_at"`
}

func (ApiCfg *ApiConfig) Refresh(w http.ResponseWriter, r *http.Request) {

	type Token struct {
		Token string `json:"token"`
	}
	tkn, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, 401, "error collecting token")
		return
	}
	requestedUser, err := ApiCfg.Db.GetUserFromToken(r.Context(), tkn)
	if err != nil {
		respondWithError(w, 401, "token doesnt exist or has expired")
		return
	}

	token, err := auth.MakeJWT(requestedUser.ID, ApiCfg.TokenSecret, time.Hour*1)

	respToken := Token{}
	respToken.Token = token

	respondWithJSON(w, 200, respToken)

}

func (ApiCfg *ApiConfig) Revoke(w http.ResponseWriter, r *http.Request) {

	type empty struct{}
	emp := empty{}

	tkn, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "error collecting token")
		return
	}

	err = ApiCfg.Db.RevokeRefreshToken(r.Context(), tkn)
	if err != nil {
		respondWithError(w, 401, "error collecting token")
		return
	}

	respondWithJSON(w, 204, emp)

}
